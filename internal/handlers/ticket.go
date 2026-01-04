package handlers

import (
	"net/http"
	"ticket-blitz/internal/repository"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	repo *repository.Repo
}

func NewTicketHandler(repo *repository.Repo) *TicketHandler {
	return &TicketHandler{repo: repo}
}

func (h *TicketHandler) ResetInventory(c *gin.Context) {
	if err := h.repo.ResetInventory(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset inventory"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inventory reset to 100"})
}

func (h *TicketHandler) BuyTicket(c *gin.Context) {
	success, err := h.repo.BuyTicketAtomic()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal system error"})
		return
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Ticket purchased!"})
	} else {
		c.JSON(http.StatusGone, gin.H{"message": "Sold out!"})
	}
}
