package bootstrap

import (
	"fmt"
	"log"

	"workshop-registration/backend/internal/adapter/storage/postgres"
	"workshop-registration/backend/internal/adapter/storage/postgres/models"
	"workshop-registration/backend/internal/adapter/storage/postgres/repository"
	"workshop-registration/backend/internal/config"
	"workshop-registration/backend/internal/core/service"

	"gorm.io/gorm"
)

type App struct {
	Config              *config.Config
	DB                  *gorm.DB
	AuthService         *service.AuthService
	WorkshopService     *service.WorkshopService
	RegistrationService *service.RegistrationService
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

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.UserModel{},
		&models.WorkshopModel{},
		&models.RegistrationModel{},
	); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := SeedAdmin(db, cfg); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("admin seeding failed: %w", err)
	}

	log.Println("Database initialized successful")

	workshopRepo := repository.NewWorkshopRepository(db)

	workshopService := service.NewWorkshopService(workshopRepo)

	registrationService := service.NewRegistrationService(db)

	return &App{
		Config:              cfg,
		DB:                  db,
		AuthService:         service.NewAuthService(db, cfg.JWTSecret),
		WorkshopService:     workshopService,
		RegistrationService: registrationService,
	}, nil
}
