package handlers

import (
	"net/http"

	"github.com/emmanuella-codes/olu/services"
	"github.com/emmanuella-codes/olu/validator"
	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	svc       *services.OTPService
	jwtSecret string
}

func NewOTPHandler(svc *services.OTPService, jwtSecret string) *OTPHandler {
	return &OTPHandler{svc: svc, jwtSecret: jwtSecret}
}

type sendOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type verifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (h *OTPHandler) Send(c *gin.Context) {
	var req sendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone is required",
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	if !validator.IsValidPhone(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid Nigerian phone number",
			"code":  "INVALID_PHONE",
		})
		return
	}

	phone := validator.ToE164(req.Phone)
	if err := h.svc.Send(c.Request.Context(), phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send OTP",
			"code":  "SMS_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"phone":   maskPhone(phone),
	})
}

func (h *OTPHandler) Verify(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone and code are required",
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	if !validator.IsValidPhone(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid Nigerian phone number",
			"code":  "INVALID_PHONE",
		})
		return
	}

	phone := validator.ToE164(req.Phone)
	token, err := h.svc.VerifyCode(c.Request.Context(), phone, req.Code, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired OTP",
			"code":  "INVALID_OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"phone": maskPhone(phone),
	})
}

func maskPhone(phone string) string {
	if len(phone) < 6 {
		return "****"
	}
	return phone[:4] + "******" + phone[len(phone)-4:]
}
