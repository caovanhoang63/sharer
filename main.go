package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"

	"sharer/internal/database"
	"sharer/internal/modules/page"
)

func main() {
	// Initialize database
	dbConfig := database.Config{
		DSN:     "./sharer.db",
		LogMode: logger.Silent, // Use Silent for production
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize layers
	pageRepo := page.NewRepository(db)
	pageService := page.NewService(pageRepo)
	pageController := page.NewController(pageService)

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	
	// Create Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", pageController.Home)
	r.GET("/pages", pageController.Index)
	r.POST("/", pageController.CreateFromForm)
	r.POST("/api/share", pageController.CreateFromAPI)
	r.GET("/shared/:slug", pageController.GetSharedContent)

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}
