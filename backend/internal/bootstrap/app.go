package bootstrap

import (
	"log"

	"workshop-registration/backend/internal/adapter/storage/postgres"
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

	return &App{
		Config: cfg,
		DB:     db,
	}, nil
}
