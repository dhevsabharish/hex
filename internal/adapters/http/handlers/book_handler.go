package handlers

import (
	"hex/internal/application/auth"
	"hex/internal/application/services"
	"hex/pkg/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type BookHandler struct {
	service     *services.BookService
	authService auth.AuthService
}

func NewBookHandler(service *services.BookService, authService auth.AuthService) *BookHandler {
	return &BookHandler{service: service, authService: authService}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var body struct {
		Title           string `json:"title" binding:"required"`
		Author          string `json:"author" binding:"required"`
		PublicationDate string `json:"publication_date" binding:"required"`
		Genre           string `json:"genre"`
		Availability    uint   `json:"availability"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, body.PublicationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication date format"})
		return
	}

	book := models.Book{
		Title:           body.Title,
		Author:          body.Author,
		PublicationDate: datatypes.Date(parsedDate),
		Genre:           body.Genre,
		Availability:    body.Availability,
	}

	token := c.GetHeader("Authorization")
	_, role, err := h.authService.Authenticate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "admin" && role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only admins and librarians can create books"})
		return
	}

	if err := h.service.CreateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) ViewAllBooks(c *gin.Context) {
	token := c.GetHeader("Authorization")
	_, _, err := h.authService.Authenticate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	books, err := h.service.ViewAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Title           string `json:"title" binding:"required"`
		Author          string `json:"author" binding:"required"`
		PublicationDate string `json:"publication_date" binding:"required"`
		Genre           string `json:"genre"`
		Availability    uint   `json:"availability"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.GetHeader("Authorization")
	_, role, err := h.authService.Authenticate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "admin" && role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only admins and librarians can update books"})
		return
	}

	existingBook, err := h.service.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingBook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	existingBook.Title = body.Title
	existingBook.Author = body.Author
	existingBook.Genre = body.Genre
	existingBook.Availability = body.Availability

	if body.PublicationDate != "" {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, body.PublicationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication date format"})
			return
		}
		existingBook.PublicationDate = datatypes.Date(parsedDate)
	}

	if err := h.service.UpdateBook(existingBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingBook)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	token := c.GetHeader("Authorization")
	_, role, err := h.authService.Authenticate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if role != "admin" && role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: only admins and librarians can delete books"})
		return
	}

	book, err := h.service.GetBookByID(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := h.service.DeleteBook(strconv.FormatUint(uint64(id), 10)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
