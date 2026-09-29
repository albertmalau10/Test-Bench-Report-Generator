package database

import (
	"log"
	"os"

	"valve_database/models"
	"valve_database/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB{
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath+"?_pragma=journal_mode(WAL)"), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal terhubung ke database:", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	if err = db.AutoMigrate(&models.Valve{}, &models.User{}, &models.Settings{}); err != nil {
		log.Fatal("gagal migrasi tabel:", err)
	}

	log.Println("Database siap!")

	return db
}
func SeedUsers(db *gorm.DB) {
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	operatorPassword := os.Getenv("OPERATOR_PASSWORD")

	if adminPassword == "" || operatorPassword == "" {
		log.Fatal("ADMIN_PASSWORD or OPERATOR_PASSWORD not set in .env")
	}

	seedUser(db, "admin", adminPassword, "admin")
	seedUser(db, "operator", operatorPassword, "operator")
}

func seedUser(db *gorm.DB, username, plainPassword, role string) {
	var existing models.User
	result := db.Where("username = ?", username).First(&existing)

	if result.Error == nil {
		// user sudah ada, skip
		return
	}

	hashed, err := utils.HashPassword(plainPassword)
	if err != nil {
		log.Printf("failed to hash password for %s: %v", username, err)
		return
	}

	newUser := models.User{
		Username: username,
		Password: hashed,
		Role:     role,
	}

	if result := db.Create(&newUser); result.Error != nil {
		log.Printf("failed to seed user %s: %v", username, result.Error)
		return
	}

	log.Printf("seeded user: %s (%s)", username, role)
}
