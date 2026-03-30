package handlers

import (
	"net/http"

	"github.com/emmanuella-codes/olu/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CandidateHandler struct {
	svc *services.CandidateService
}

func NewCandidateHandler(svc *services.CandidateService) *CandidateHandler {
	return &CandidateHandler{svc: svc}
}

func (h *CandidateHandler) List(c *gin.Context) {
	candidates, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  candidates,
		"count": len(candidates),
	})
}

func (h *CandidateHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_ID",
		})
		return
	}

	candidate, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}
	if candidate == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "candidate not found",
			"code":  "NOT_FOUND",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": candidate,
	})
}
