package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AdminClaims struct {
	AdminID uuid.UUID `json:"admin_id"`
	Email   string    `json:"email"`
	jwt.RegisteredClaims
}

const AdminTokenTTL = 8 * time.Hour

func RequireAdminToken(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "admin token required",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			return
		}

		claims, ok := token.Claims.(*AdminClaims)
		if !ok || claims.AdminID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			return
		}
		c.Set("admin_id", claims.AdminID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func NewAdminClaims(adminID uuid.UUID, email string, now time.Time) AdminClaims {
	return AdminClaims{
		AdminID: adminID,
		Email:   email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AdminTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "olu-admin",
			Subject:   adminID.String(),
		},
	}
}
