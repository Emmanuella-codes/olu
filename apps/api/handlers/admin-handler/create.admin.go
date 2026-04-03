package adminhandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type createAdminRequest struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

const minAdminPasswordLength = 12

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req createAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	req.Email = normalizeAdminEmail(req.Email)
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	if len(req.Password) < minAdminPasswordLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 12 characters"})
		return
	}

	existing, err := admin.AdminRepo.GetAdminByEmail(c.Request.Context(), req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("failed to check existing admin")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "an admin with this email already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash admin password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin"})
		return
	}

	newAdmin, err := admin.AdminRepo.CreateAdmin(c.Request.Context(), req.Email, string(hash))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "an admin with this email already exists"})
			return
		}
		log.Error().Err(err).Str("email", req.Email).Msg("failed to insert admin")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin"})
		return
	}

	log.Info().Str("email", newAdmin.Email).Str("id", newAdmin.ID.String()).Msg("admin created")
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":    newAdmin.ID,
			"email": newAdmin.Email,
		},
	})
}

func normalizeAdminEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
