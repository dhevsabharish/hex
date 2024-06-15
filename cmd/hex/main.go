package main

import (
	"hex/config"
	"hex/internal/adapters/http/handlers"
	"hex/internal/adapters/persistence"
	"hex/internal/application/services"

	"hex/internal/adapters/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	authService := auth.NewRailsAuthService("http://localhost:3000")
	bookRepo := persistence.NewBookRepository(config.DB)
	bookService := services.NewBookService(*bookRepo, authService)
	bookHandler := handlers.NewBookHandler(bookService)

	borrowingRepo := persistence.NewBorrowingRepository(config.DB)
	borrowingService := services.NewBorrowingService(*bookRepo, *borrowingRepo, authService)
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

	r.Run(":8080")
}
