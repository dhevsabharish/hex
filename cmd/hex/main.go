package main

import (
	"hex/config"
	"hex/internal/adapters/auth"
	"hex/internal/adapters/http/handlers"
	"hex/internal/adapters/persistence"
	"hex/internal/application/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()

	authService := auth.NewRailsAuthService(cfg.RailsAPIURL)
	bookRepo := persistence.NewBookRepository(cfg.DB)
	bookService := services.NewBookService(*bookRepo, authService)
	bookHandler := handlers.NewBookHandler(bookService)

	borrowingRepo := persistence.NewBorrowingRepository(cfg.DB)
	borrowingService := services.NewBorrowingService(*bookRepo, *borrowingRepo)
	borrowingHandler := handlers.NewBorrowingHandler(borrowingService, authService)

	r := gin.Default()

	r.POST("/books", bookHandler.CreateBook)
	r.GET("/books", bookHandler.ViewAllBooks)
	r.PUT("/books/:id", bookHandler.UpdateBook)
	r.DELETE("/books/:id", bookHandler.DeleteBook)

	r.POST("/borrow", borrowingHandler.BorrowBook)
	r.POST("/return", borrowingHandler.ReturnBook)
	r.GET("/my-borrowings", borrowingHandler.GetMyBorrowings)
	r.GET("/borrowing-records", borrowingHandler.GetAllBorrowingRecords)

	r.Run(":" + cfg.Port)
}
