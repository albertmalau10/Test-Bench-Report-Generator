package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"valve_database/controllers"
	"valve_database/database"
	"valve_database/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: .env file not found, using system environment defaults")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	db := database.ConnectDB()
	database.SeedUsers(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to obtain SQL DB handle:", err)
	}
	defer sqlDB.Close()

	_ = os.MkdirAll("./images", os.ModePerm)
	_ = os.MkdirAll("./datasheets", os.ModePerm)

	router := gin.Default()

	// CORS configuration for local dev and local network access
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static asset routing
	router.Static("/images", "./images")
	router.Static("/datasheets", "./datasheets")

	api := router.Group("/api")
	{
		// Health & Auth
		api.GET("/ping", controllers.Ping)
		api.POST("/login", controllers.Login(db))

		// Valve Master Data
		api.GET("/valves", controllers.GetValves(db))
		api.GET("/valves/:id", controllers.GetValveByID(db))
		api.PUT("/valves/:id", middleware.AuthRequired(), controllers.UpdateValve(db))
		api.POST("/valves", middleware.AuthRequired(), middleware.AdminOnly(), controllers.CreateValve(db))
		api.DELETE("/valves/:id", middleware.AuthRequired(), middleware.AdminOnly(), controllers.DeleteValve(db))
		api.POST("/valves/:id/image", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveImage(db))
		api.POST("/valves/:id/datasheet", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveDatasheet(db))

		// Test Records & Telemetry Logs
		api.GET("/records", middleware.AuthRequired(), controllers.GetTestRecords(db))
		api.POST("/records", middleware.AuthRequired(), controllers.SaveTestRecord(db))

		// System Settings
		api.GET("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.GetSettings(db))
		api.PUT("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UpdateSettings(db))

		// OPC UA Telemetry & Actuation
		api.GET("/opcua/data", middleware.AuthRequired(), controllers.GetOpcData(db))
		api.GET("/opcua/status", middleware.AuthRequired(), controllers.GetOpcStatus(db))
		api.POST("/ctrlx/start", middleware.AuthRequired(), controllers.StartOutput(db))
		api.POST("/ctrlx/stop", middleware.AuthRequired(), controllers.StopOutput(db))
	}

	if err := router.Run(":" + appPort); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}