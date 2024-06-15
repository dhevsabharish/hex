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

	userID, role, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "member" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: only members can borrow books"})
		return
	}

	memberID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BorrowBook(body.BookID, uint(memberID)); err != nil {
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

	userID, role, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "member" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: only members can return books"})
		return
	}

	memberID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ReturnBook(body.BorrowingRecordID, uint(memberID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func (h *BorrowingHandler) GetMyBorrowings(c *gin.Context) {
	userID, role, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "member" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: only members can view their borrowings"})
		return
	}

	memberID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	borrowingRecords, err := h.service.GetMyBorrowings(uint(memberID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrowingRecords)
}

func (h *BorrowingHandler) GetAllBorrowingRecords(c *gin.Context) {
	_, role, err := h.authService.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "admin" && role != "librarian" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: only admins and librarians can view all borrowing records"})
		return
	}

	borrowingRecords, err := h.service.GetAllBorrowingRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrowingRecords)
}