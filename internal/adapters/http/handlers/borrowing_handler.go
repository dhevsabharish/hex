package handlers

import (
	"net/http"
	"strconv"

	"hex/internal/application/auth"
	"hex/internal/application/services"

	"github.com/gin-gonic/gin"
)

type BorrowingHandler struct {
	service     services.BorrowingService
	authService auth.AuthService
}

func NewBorrowingHandler(service services.BorrowingService, authService auth.AuthService) *BorrowingHandler {
	return &BorrowingHandler{service: service, authService: authService}
}

func (h *BorrowingHandler) BorrowBook(c *gin.Context) {
	var body struct {
		BookID uint `json:"book_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	memberID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BorrowBook(body.BookID, uint(memberID), c.GetHeader("Authorization")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book borrowed successfully"})
}

func (h *BorrowingHandler) ReturnBook(c *gin.Context) {
	var body struct {
		BorrowingRecordID uint `json:"borrowing_record_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	memberID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ReturnBook(body.BorrowingRecordID, uint(memberID), c.GetHeader("Authorization")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func (h *BorrowingHandler) GetMyBorrowings(c *gin.Context) {
	userID, _, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	memberID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	borrowingRecords, err := h.service.GetMyBorrowings(uint(memberID), c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrowingRecords)
}

func (h *BorrowingHandler) GetAllBorrowingRecords(c *gin.Context) {
	borrowingRecords, err := h.service.GetAllBorrowingRecords(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrowingRecords)
}
