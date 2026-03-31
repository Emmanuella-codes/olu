package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AllMessages(c *gin.Context) {
	msgs := h.store.All()
	// Reverse so newest is first
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": msgs,
		"count":    len(msgs),
		"stats":    h.store.Stats(),
	})
}

func (h *Handler) ByPhone(c *gin.Context) {
	phone := c.Param("phone")
	msgs := h.store.ByPhone(phone)
	// reverse so newest is first
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	c.JSON(http.StatusOK, gin.H{
		"phone":    phone,
		"messages": msgs,
		"count":    len(msgs),
	})
}

func (h *Handler) Clear(c *gin.Context) {
	h.store.Clear()
	c.JSON(http.StatusOK, gin.H{"message": "all messages cleared"})
}

func (h *Handler) Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"stats": h.store.Stats()})
}
