package database

import (
	"log"
	"os"
	"time"

	"valve_database/models"
	"valve_database/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data.db"
	}

	dsn := dbPath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(1 * time.Hour)
	}

	// AutoMigrate now includes TestRecord
	if err = db.AutoMigrate(&models.Valve{}, &models.User{}, &models.Settings{}, &models.TestRecord{}); err != nil {
		log.Fatalf("Table migration failed: %v", err)
	}

	initDefaultSettings(db)

	log.Println("Database connection and schema ready.")
	return db
}

func initDefaultSettings(db *gorm.DB) {
	var s models.Settings
	if err := db.First(&s, 1).Error; err != nil {
		defaultSettings := models.Settings{
			ID:                  1,
			OpcUaAddress:        "opc.tcp://127.0.0.1:4840",
			OpcNodePressure:     "ns=2;s=testbench/sensors/pressure",
			OpcNodeCommand:      "ns=2;s=testbench/sensors/command",
			OpcNodeFeedback:     "ns=2;s=testbench/sensors/feedback",
			OpcNodeFlow:         "ns=2;s=testbench/sensors/flow",
			OpcNodeOutput:       "ns=2;s=plc/app/Application/sym/PLC_PRG/output",
			OpcUaSecurityPolicy: "None",
			OpcUaSecurityMode:   "None",
		}
		db.Create(&defaultSettings)
	}
}

func SeedUsers(db *gorm.DB) {
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	operatorPassword := os.Getenv("OPERATOR_PASSWORD")

	if adminPassword == "" {
		adminPassword = "adminpassword"
	}
	if operatorPassword == "" {
		operatorPassword = "operatorpassword"
	}

	seedUser(db, "admin", adminPassword, "admin")
	seedUser(db, "operator", operatorPassword, "operator")
}

func seedUser(db *gorm.DB, username, plainPassword, role string) {
	var existing models.User
	if err := db.Where("username = ?", username).First(&existing).Error; err == nil {
		return
	}

	hashed, err := utils.HashPassword(plainPassword)
	if err != nil {
		return
	}

	newUser := models.User{
		Username: username,
		Password: hashed,
		Role:     role,
	}
	db.Create(&newUser)
}