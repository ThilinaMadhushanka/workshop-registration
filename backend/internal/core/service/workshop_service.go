package service

import (
	"errors"
	"strings"
	"time"

	"workshop-registration/backend/internal/adapter/http/dto"
	"workshop-registration/backend/internal/adapter/storage/postgres/models"
	"workshop-registration/backend/internal/adapter/storage/postgres/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrWorkshopInvalid = errors.New("invalid workshop details")
	ErrCapacityTooLow  = errors.New("capacity is below active registrations")
	ErrWorkshopExists  = errors.New("workshop code already exists")
)

type WorkshopService struct {
	Repo *repository.WorkshopRepository
}

func NewWorkshopService(
	repo *repository.WorkshopRepository,
) *WorkshopService {
	return &WorkshopService{Repo: repo}
}

func validateWorkshop(req dto.WorkshopRequest) (time.Time, error) {
	if strings.TrimSpace(req.Code) == "" ||
		strings.TrimSpace(req.Title) == "" ||
		strings.TrimSpace(req.Instructor) == "" ||
		req.Capacity <= 0 {
		return time.Time{}, ErrWorkshopInvalid
	}

	if len(req.Code) > 50 ||
		len(req.Title) > 200 ||
		len(req.Instructor) > 150 {
		return time.Time{}, ErrWorkshopInvalid
	}

	switch req.Status {
	case "scheduled", "cancelled", "completed":
	default:
		return time.Time{}, ErrWorkshopInvalid
	}

	start, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return time.Time{}, errors.New(
			"startsAt must use RFC3339 date-time format",
		)
	}

	return start, nil
}

func (s *WorkshopService) Create(
	req dto.WorkshopRequest,
) (*models.WorkshopModel, error) {

	start, err := validateWorkshop(req)
	if err != nil {
		return nil, err
	}

	workshop := models.WorkshopModel{
		Code:       strings.TrimSpace(req.Code),
		Title:      strings.TrimSpace(req.Title),
		Instructor: strings.TrimSpace(req.Instructor),
		StartsAt:   start,
		Capacity:   req.Capacity,
		Status:     req.Status,
	}

	if err := s.Repo.Create(&workshop); err != nil {
		return nil, err
	}

	return &workshop, nil
}

func (s *WorkshopService) GetAll() (
	[]models.WorkshopModel, error,
) {
	return s.Repo.GetAll()
}

func (s *WorkshopService) GetByID(
	id uint,
) (*models.WorkshopModel, error) {
	return s.Repo.GetByID(id)
}

func (s *WorkshopService) Update(
	id uint,
	req dto.WorkshopRequest,
) (*models.WorkshopModel, error) {

	start, err := validateWorkshop(req)
	if err != nil {
		return nil, err
	}

	var workshop models.WorkshopModel

	err = s.Repo.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&workshop, id).Error; err != nil {
			return err
		}

		var activeCount int64

		if err := tx.Model(&models.RegistrationModel{}).
			Where("workshop_id = ? AND status = ?", id, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int64(req.Capacity) < activeCount {
			return ErrCapacityTooLow
		}

		workshop.Code = strings.TrimSpace(req.Code)
		workshop.Title = strings.TrimSpace(req.Title)
		workshop.Instructor = strings.TrimSpace(req.Instructor)
		workshop.StartsAt = start
		workshop.Capacity = req.Capacity
		workshop.Status = req.Status

		return tx.Save(&workshop).Error
	})

	if err != nil {
		return nil, err
	}

	return &workshop, nil
}

var ErrInvalidFilter = errors.New(
	"invalid workshop filter",
)

func (s *WorkshopService) GetFiltered(
	from string,
	to string,
	status string,
	availableOnly bool,
) ([]models.WorkshopModel, error) {

	filter := repository.WorkshopFilter{
		AvailableOnly: availableOnly,
	}

	if status != "" {
		switch status {
		case "scheduled", "cancelled", "completed":
			filter.Status = status
		default:
			return nil, ErrInvalidFilter
		}
	}

	if from != "" {
		parsed, err := time.Parse("2006-01-02", from)
		if err != nil {
			return nil, ErrInvalidFilter
		}
		filter.From = &parsed
	}

	if to != "" {
		parsed, err := time.Parse("2006-01-02", to)
		if err != nil {
			return nil, ErrInvalidFilter
		}

		exclusive := parsed.AddDate(0, 0, 1)
		filter.ToExclusive = &exclusive
	}

	if filter.From != nil && filter.ToExclusive != nil {
		if !filter.From.Before(*filter.ToExclusive) {
			return nil, ErrInvalidFilter
		}
	}

	return s.Repo.GetFiltered(filter)
}
