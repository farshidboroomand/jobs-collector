package application

import (
	"net/http"
	"strconv"

	"github.com/farshidboroomand/jobs-collector/internal/domain"
	"github.com/farshidboroomand/jobs-collector/internal/service"
	"github.com/gin-gonic/gin"
)

// ContinentHandler handles HTTP requests related to continents, utilizing the Continent service for business logic.
type ContinentHandler struct {
	services *service.Services
}

// NewContinentHandler creates a new instance of ContinentHandler with the provided services.
func NewContinentHandler(services *service.Services) *ContinentHandler {
	return &ContinentHandler{
		services: services,
	}
}

// FetchContinentByID retrieves a continent by its ID from the Continent service and returns it as a JSON response.
func (h *ContinentHandler) FetchContinentByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid continent id"})
		return
	}

	continent, err := h.services.Continent.GetContinentByID(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, continent)
}

// CreateContinent handles the creation of a new continent based on the incoming JSON request.
func (h *ContinentHandler) CreateContinent(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	continent := &domain.Continent{
		Name: req.Name,
	}

	if err := h.services.Continent.PublishContinent(c.Request.Context(), continent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, continent)
}
