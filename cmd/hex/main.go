package main

import (
	"hex/config"
	"hex/internal/adapters/http/handlers"
	"hex/internal/adapters/persistence"
	"hex/internal/application/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	bookRepo := persistence.NewBookRepository(config.DB)
	bookService := services.NewBookService(*bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	borrowingRepo := persistence.NewBorrowingRepository(config.DB)
	borrowingService := services.NewBorrowingService(*borrowingRepo)
	borrowingHandler := handlers.NewBorrowingHandler(borrowingService)

	r := gin.Default()

	r.POST("/books", bookHandler.CreateBook)
	r.GET("/books", bookHandler.ViewAllBooks)
	r.PUT("/books/:id", bookHandler.UpdateBook)
	r.DELETE("/books/:id", bookHandler.DeleteBook)

	r.POST("/return/:id", borrowingHandler.ReturnBook)
	r.GET("/records", borrowingHandler.ShowAllRecords)
	r.GET("/records/user/:userID", borrowingHandler.ShowRecordsByUserID)

	r.Run(":8080")
}
