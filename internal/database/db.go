package database

import (
	"log"
	"os"

	"github.com/radhikadarode/finance-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "finance.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate all models
	err = DB.AutoMigrate(&models.User{}, &models.FinancialRecord{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database initialized successfully")
	seedAdminUser()
}

// seedAdminUser creates a default admin user if none exists
func seedAdminUser() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin@123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash seed password: %v", err)
		return
	}

	admin := models.User{
		Name:     "System Admin",
		Email:    "admin@finance.local",
		Password: string(hash),
		Role:     models.RoleAdmin,
		IsActive: true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin user: %v", err)
		return
	}

	log.Println("Default admin user created: admin@finance.local / admin@123")
}
