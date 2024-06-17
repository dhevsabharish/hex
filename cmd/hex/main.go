package main

import (
	"hex/config"
	"hex/internal/adapters/auth"
	"hex/internal/adapters/cors"
	"hex/internal/adapters/http/handlers"
	"hex/internal/adapters/persistence"
	"hex/internal/adapters/seeder"
	"hex/internal/application/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the configuration
	cfg := config.NewConfig()

	// Run seeding if the environment variable is set to true
	if cfg.SeedDatabase {
		if err := seeder.Seed(cfg.DB); err != nil {
			log.Fatalf("Error seeding database: %v", err)
		}
		log.Println("Database seeding completed.")
	} else {
		log.Println("Database seeding is disabled.")
	}

	// Initialize authentication service
	authService := auth.NewRailsAuthService(cfg.RailsAPIURL)

	// Initialize repositories
	bookRepo := persistence.NewBookRepository(cfg.DB)
	borrowingRepo := persistence.NewBorrowingRepository(cfg.DB)

	// Initialize services
	bookService := services.NewBookService(*bookRepo, cfg.Logger)
	borrowingService := services.NewBorrowingService(*bookRepo, *borrowingRepo, cfg.Logger)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService, authService)
	borrowingHandler := handlers.NewBorrowingHandler(borrowingService, authService)

	// Setup Gin router
	r := gin.Default()

	// Configure CORS
	cors.ConfigureCORS(r)

	// Define routes
	r.POST("/books", bookHandler.CreateBook)
	r.GET("/books", bookHandler.ViewAllBooks)
	r.PUT("/books/:id", bookHandler.UpdateBook)
	r.DELETE("/books/:id", bookHandler.DeleteBook)

	r.POST("/borrow", borrowingHandler.BorrowBook)
	r.POST("/return", borrowingHandler.ReturnBook)
	r.GET("/my-borrowings", borrowingHandler.GetMyBorrowings)
	r.GET("/borrowing-records", borrowingHandler.GetAllBorrowingRecords)

	// Run the server
	r.Run(":" + cfg.Port)
}
