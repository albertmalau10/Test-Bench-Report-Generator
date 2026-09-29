package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"valve_database/controllers"
	"valve_database/database"
	"valve_database/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env tidak ditemukan, pakai default")
	}
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	db := database.ConnectDB()
	database.SeedUsers(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("gagal mengambil koneksi sql:", err)
	}
	defer sqlDB.Close()

	// FIX
	os.MkdirAll("./images", os.ModePerm)
	os.MkdirAll("./datasheets", os.ModePerm)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/ping", controllers.Ping)
		api.POST("/login", controllers.Login(db))

		api.GET("/valves", controllers.GetValves(db))
		api.GET("/valves/:id", controllers.GetValveByID(db))

		api.PUT("/valves/:id", middleware.AuthRequired(), controllers.UpdateValve(db))

		api.POST("/valves", middleware.AuthRequired(), middleware.AdminOnly(), controllers.CreateValve(db))
		api.DELETE("/valves/:id", middleware.AuthRequired(), middleware.AdminOnly(), controllers.DeleteValve(db))
		api.POST("/valves/:id/image", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveImage(db))
		api.POST("/valves/:id/datasheet", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveDatasheet(db))

		api.GET("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.GetSettings(db))
		api.PUT("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UpdateSettings(db))

		// OPC UA routes
		api.GET("/opcua/data", middleware.AuthRequired(), controllers.GetOpcData(db))
		api.POST("/ctrlx/start", middleware.AuthRequired(), controllers.StartOutput(db))
		api.GET("/opcua/status", middleware.AuthRequired(), controllers.GetOpcStatus(db))
		api.POST("/ctrlx/stop", middleware.AuthRequired(), controllers.StopOutput(db))

		
	}
	router.Static("/images", "./images")
	router.Static("/datasheets", "./datasheets")

	router.Run(":" + appPort)
}
