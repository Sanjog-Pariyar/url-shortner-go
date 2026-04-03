package handlers

import (
	"errors"
	"net/http"

	"github.com/sanjog-pariyar/url-shorten-api/internal/models"
	"github.com/sanjog-pariyar/url-shorten-api/internal/services"
	"github.com/sanjog-pariyar/url-shorten-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service services.URLService
}

func NewURLHandler(service services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID *uint64
	if id, exists := c.Get("userID"); exists {
		uid := id.(uint64)
		userID = &uid
	}

	resp, err := h.service.CreateShortURL(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		if errors.Is(err, services.ErrURLExpired) {
			c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func (h *URLHandler) GetStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	stats, err := h.service.GetStats(shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var userID *uint64
	if id, exists := c.Get("userID"); exists {
		uid := id.(uint64)
		userID = &uid
	}

	err := h.service.DeleteURL(shortCode, userID)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this URL"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}