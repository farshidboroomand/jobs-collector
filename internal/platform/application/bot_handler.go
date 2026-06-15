package application

import (
	"net/http"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/farshidboroomand/jobs-collector/internal/service"
	"github.com/gin-gonic/gin"
)

// BotHandler handles HTTP requests related to bots, utilizing the Bot service for business logic.
type BotHandler struct {
	services *service.Services
}

// NewBotHandler creates a new instance of BotHandler with the provided services.
func NewBotHandler(services *service.Services) *BotHandler {
	return &BotHandler{
		services: services,
	}
}

// List retrieves a list of bots from the Bot service and returns them as a JSON response.
func (h *BotHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	bots, _, err := h.services.Bot.ListBots(ctx, 1, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bots)
}

// CreateBot handles the creation of a new bot based on the incoming JSON request.
func (h *BotHandler) CreateBot(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		JobPosition string `json:"job_position" binding:"required"`
		CountryID   int    `json:"country_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	bot := &domain.Bot{
		Title:       req.Title,
		JobPosition: req.JobPosition,
		CountryID:   req.CountryID,
		Status:      1,
		IsActive:    true,
	}

	if err := h.services.Bot.PublishBot(c.Request.Context(), bot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bot)
}
