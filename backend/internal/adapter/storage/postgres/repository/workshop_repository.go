package repository

import (
	"gorm.io/gorm"
	"time"
	"workshop-registration/backend/internal/adapter/storage/postgres/models"
)

type WorkshopRepository struct {
	DB *gorm.DB
}

func NewWorkshopRepository(db *gorm.DB) *WorkshopRepository {
	return &WorkshopRepository{DB: db}
}

func (r *WorkshopRepository) Create(
	workshop *models.WorkshopModel,
) error {
	return r.DB.Create(workshop).Error
}

func (r *WorkshopRepository) GetAll() (
	[]models.WorkshopModel, error,
) {
	var workshops []models.WorkshopModel

	err := r.DB.Order("starts_at ASC").
		Find(&workshops).Error

	return workshops, err
}

func (r *WorkshopRepository) GetByID(
	id uint,
) (*models.WorkshopModel, error) {

	var workshop models.WorkshopModel

	err := r.DB.First(&workshop, id).Error
	if err != nil {
		return nil, err
	}

	return &workshop, nil
}

func (r *WorkshopRepository) ActiveCount(
	workshopID uint,
) (int64, error) {

	var count int64

	err := r.DB.Model(&models.RegistrationModel{}).
		Where(
			"workshop_id = ? AND status = ?",
			workshopID,
			"active",
		).
		Count(&count).Error

	return count, err
}

type WorkshopFilter struct {
	From          *time.Time
	ToExclusive   *time.Time
	Status        string
	AvailableOnly bool
}

func (r *WorkshopRepository) GetFiltered(
	filter WorkshopFilter,
) ([]models.WorkshopModel, error) {

	var workshops []models.WorkshopModel

	query := r.DB.Model(&models.WorkshopModel{})

	if filter.From != nil {
		query = query.Where(
			"starts_at >= ?", *filter.From,
		)
	}

	if filter.ToExclusive != nil {
		query = query.Where(
			"starts_at < ?", *filter.ToExclusive,
		)
	}

	if filter.Status != "" {
		query = query.Where(
			"status = ?", filter.Status,
		)
	}

	if filter.AvailableOnly {
		query = query.Where(`
            (
                SELECT COUNT(*)
                FROM registrations
                WHERE registrations.workshop_id = workshops.id
                AND registrations.status = 'active'
            ) < workshops.capacity
        `)
	}

	err := query.Order("starts_at ASC").
		Find(&workshops).Error

	return workshops, err
}
