package service

import (
	"errors"
	"strings"
	"time"

	"workshop-registration/backend/internal/adapter/storage/postgres/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrWorkshopFull     = errors.New("workshop is full")
	ErrWorkshopClosed   = errors.New("workshop is not open for registration")
	ErrInvalidAttendee  = errors.New("invalid attendee details")
	ErrAlreadyCancelled = errors.New("registration already cancelled")
)

type RegistrationService struct {
	DB *gorm.DB
}

func NewRegistrationService(db *gorm.DB) *RegistrationService {
	return &RegistrationService{DB: db}
}

func (s *RegistrationService) Create(
	workshopID uint,
	name string,
	email string,
	staffID uint,
) (*models.RegistrationModel, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if workshopID == 0 || staffID == 0 ||
		name == "" || email == "" ||
		len(name) > 150 || len(email) > 255 {
		return nil, ErrInvalidAttendee
	}

	var registration models.RegistrationModel

	err := s.DB.Transaction(func(tx *gorm.DB) error {

		var workshop models.WorkshopModel

		if err := tx.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&workshop, workshopID).Error; err != nil {
			return err
		}

		if workshop.Status != "scheduled" {
			return ErrWorkshopClosed
		}

		var activeCount int64

		if err := tx.Model(&models.RegistrationModel{}).
			Where(
				"workshop_id = ? AND status = ?",
				workshopID,
				"active",
			).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if activeCount >= int64(workshop.Capacity) {
			return ErrWorkshopFull
		}

		registration = models.RegistrationModel{
			WorkshopID:    workshopID,
			AttendeeName:  name,
			AttendeeEmail: email,
			Status:        "active",
			RegisteredBy:  staffID,
			RegisteredAt:  time.Now().UTC(),
		}

		return tx.Create(&registration).Error
	})

	if err != nil {
		return nil, err
	}

	return &registration, nil
}
func (s *RegistrationService) GetByWorkshop(
	workshopID uint,
) ([]models.RegistrationModel, error) {

	var workshop models.WorkshopModel

	if err := s.DB.First(&workshop, workshopID).Error; err != nil {
		return nil, err
	}

	var registrations []models.RegistrationModel

	err := s.DB.
		Where("workshop_id = ?", workshopID).
		Order("registered_at DESC, id DESC").
		Find(&registrations).Error

	return registrations, err
}
func (s *RegistrationService) Cancel(
	registrationID uint,
	staffID uint,
) (*models.RegistrationModel, error) {

	if registrationID == 0 || staffID == 0 {
		return nil, ErrInvalidAttendee
	}

	var registration models.RegistrationModel

	err := s.DB.Transaction(func(tx *gorm.DB) error {

		var lookup models.RegistrationModel

		if err := tx.First(&lookup, registrationID).Error; err != nil {
			return err
		}

		var workshop models.WorkshopModel

		if err := tx.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&workshop, lookup.WorkshopID).Error; err != nil {
			return err
		}

		if err := tx.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&registration, registrationID).Error; err != nil {
			return err
		}

		if registration.Status == "cancelled" {
			return ErrAlreadyCancelled
		}

		now := time.Now().UTC()

		registration.Status = "cancelled"
		registration.CancelledBy = &staffID
		registration.CancelledAt = &now

		return tx.Save(&registration).Error
	})

	if err != nil {
		return nil, err
	}

	return &registration, nil
}
