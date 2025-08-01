package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"

	"sharer/internal/database"
	"sharer/internal/modules/category"
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

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryController := category.NewController(categoryService)

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.Default()

	// Routes
	r.GET("/", pageController.Home)
	r.GET("/pages", pageController.Index)
	r.POST("/", pageController.CreateFromForm)
	r.POST("/api/share", pageController.CreateFromAPI)
	r.GET("/shared/:slug", pageController.GetSharedContent)

	// Category routes
	r.GET("/categories", categoryController.Index)
	r.GET("/categories/create", categoryController.Create)
	r.POST("/categories", categoryController.Store)
	r.GET("/categories/:id", categoryController.Show)
	r.GET("/categories/:id/edit", categoryController.Edit)
	r.PUT("/categories/:id", categoryController.Update)
	r.DELETE("/categories/:id", categoryController.Delete)
	r.GET("/api/categories", categoryController.GetAllForDropdown)

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}
