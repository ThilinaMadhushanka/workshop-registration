package bootstrap

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"workshop-registration/backend/internal/adapter/storage/postgres/models"
	"workshop-registration/backend/internal/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB, cfg *config.Config) error {
	email := strings.ToLower(strings.TrimSpace(cfg.AdminEmail))

	var existing models.UserModel

	err := db.Where("email = ?", email).First(&existing).Error

	if err == nil {
		if existing.Role != "admin" {
			return fmt.Errorf("seed email already belongs to a non-admin user")
		}
		log.Println("Admin account already exists")
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.AdminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	admin := models.UserModel{
		Name:         cfg.AdminName,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Initial Admin account created")
	return nil
}
