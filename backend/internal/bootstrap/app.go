package bootstrap

import (
	"fmt"
	"log"

	"workshop-registration/backend/internal/adapter/storage/postgres"
	"workshop-registration/backend/internal/adapter/storage/postgres/models"
	"workshop-registration/backend/internal/config"

	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
}

func NewApp() (*App, error) {

	cfg, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}

	db, err := postgres.ConnectDB(cfg)

	if err != nil {
		return nil, err
	}

	log.Println("Database connection successful")

	err = db.AutoMigrate(
		&models.UserModel{},
		&models.WorkshopModel{},
		&models.RegistrationModel{},
	)

	if err != nil {
		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			_ = sqlDB.Close()
		}
		return nil, fmt.Errorf(
			"migration failed: %w", err,
		)
	}

	return &App{
		Config: cfg,
		DB:     db,
	}, nil
}
