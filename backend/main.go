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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/images", "./images")
	router.Static("/datasheets", "./datasheets")

	router.GET("/ping", controllers.Ping)
	router.POST("/login", controllers.Login(db))

	router.GET("/valves", controllers.GetValves(db))
	router.GET("/valves/:id", controllers.GetValveByID(db))

	router.PUT("/valves/:id", middleware.AuthRequired(), controllers.UpdateValve(db))

	router.POST("/valves", middleware.AuthRequired(), middleware.AdminOnly(), controllers.CreateValve(db))
	router.DELETE("/valves/:id", middleware.AuthRequired(), middleware.AdminOnly(), controllers.DeleteValve(db))
	router.POST("/valves/:id/image", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveImage(db))
	router.POST("/valves/:id/datasheet", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UploadValveDatasheet(db))

	router.GET("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.GetSettings(db))
	router.PUT("/settings", middleware.AuthRequired(), middleware.AdminOnly(), controllers.UpdateSettings(db))
	
	//OPC UA Routes
	router.GET("/opcua/data", middleware.AuthRequired(), controllers.GetOpcData(db))
	router.POST("/ctrlx/start", middleware.AuthRequired(), controllers.StartOutput(db))
	router.GET("/opcua/status", middleware.AuthRequired(), controllers.GetOpcStatus(db))
	router.POST("/ctrlx/stop", middleware.AuthRequired(), controllers.StopOutput(db))

	router.Run(":" + appPort)
}
