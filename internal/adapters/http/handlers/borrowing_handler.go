package handlers

import (
	"hex/internal/application/services"
	"hex/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BorrowingHandler struct {
	service *services.BorrowingService
}

func NewBorrowingHandler(service *services.BorrowingService) *BorrowingHandler {
	return &BorrowingHandler{service: service}
}

func (h *BorrowingHandler) BorrowBook(c *gin.Context) {
	var record models.BorrowingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BorrowBook(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *BorrowingHandler) ReturnBook(c *gin.Context) {
	recordID := c.Param("id")
	if err := h.service.ReturnBook(recordID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func (h *BorrowingHandler) ShowAllRecords(c *gin.Context) {
	records, err := h.service.ShowAllRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *BorrowingHandler) ShowRecordsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	records, err := h.service.ShowRecordsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}
