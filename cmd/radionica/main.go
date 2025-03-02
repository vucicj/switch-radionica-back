package main

import (
	"database/sql"
	"fmt"
	"log"

	"blazperic/radionica/config"
	"blazperic/radionica/internal/api"

	_ "blazperic/radionica/docs" // Import generated docs

	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"     // Swagger UI files
	ginSwagger "github.com/swaggo/gin-swagger" // Gin middleware for Swagger
)

// @title Radionica API
// @version 1.0
// @description This is the API for the Radionica application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	server := api.NewServer(db, cfg)
	router := api.SetupRouter(server)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
