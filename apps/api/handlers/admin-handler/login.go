package adminhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/emmanuella-codes/olu/middleware"
	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	jwtSecret string
}

// Keeps the missing-user path doing real bcrypt work to reduce timing leaks.
const dummyPasswordHash = "$2a$10$8V7sQ8bE6l2mQJj8m9xP0.3wG8m4r2QY2k6c1YVQy5m2Q3lRr9s7K"

func NewAdminHandler(jwtSecret string) *AdminHandler {
	return &AdminHandler{jwtSecret: jwtSecret}
}

type adminLoginRequest struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}
	req.Email = normalizeAdminEmail(req.Email)
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}

	adminU, err := admin.AdminRepo.GetAdminByEmail(c.Request.Context(), req.Email)
	if err != nil || adminU == nil || !adminU.IsActive {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(req.Password))
		log.Warn().Str("email", req.Email).Str("ip", c.ClientIP()).Msg("admin login failed: user not found or inactive")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminU.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn().Str("email", req.Email).Str("ip", c.ClientIP()).Msg("admin login failed: invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := h.issueAdminToken(adminU.ID, adminU.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
		return
	}

	go func(adminID uuid.UUID) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := admin.AdminRepo.UpdateAdminLastLogin(ctx, adminID); err != nil {
			log.Warn().Err(err).Str("admin_id", adminID.String()).Msg("failed to update admin last login")
		}
	}(adminU.ID)
	log.Info().Str("email", adminU.Email).Msg("admin logged in successfully")
	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": int(middleware.AdminTokenTTL.Seconds()),
	})
}

func (h *AdminHandler) issueAdminToken(adminID uuid.UUID, email string) (string, error) {
	claims := middleware.NewAdminClaims(adminID, email, time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
