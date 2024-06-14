package handlers

import (
	"hex/internal/application/services"
	"hex/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{service: service}
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

	book := models.Book{
		Title:        body.Title,
		Author:       body.Author,
		Genre:        body.Genre,
		Availability: body.Availability,
	}

	token := c.GetHeader("Authorization")
	if err := h.service.CreateBook(&book, body.PublicationDate, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) ViewAllBooks(c *gin.Context) {
	token := c.GetHeader("Authorization")
	books, err := h.service.ViewAllBooks(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
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

	token := c.GetHeader("Authorization")
	if err := h.service.UpdateBook(existingBook, body.PublicationDate, token); err != nil {
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

	book, err := h.service.GetBookByID(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	token := c.GetHeader("Authorization")
	if err := h.service.DeleteBook(strconv.FormatUint(uint64(id), 10), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
