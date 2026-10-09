package api

import (
	"errors"
	"net/http"
	"time"

	"workshop-registration/backend/internal/adapter/http/dto"
	"workshop-registration/backend/internal/adapter/storage/postgres/models"
	"workshop-registration/backend/internal/core/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegistrationHandler struct {
	Service *service.RegistrationService
}

func NewRegistrationHandler(
	s *service.RegistrationService,
) *RegistrationHandler {
	return &RegistrationHandler{Service: s}
}

func registrationResponse(
	r models.RegistrationModel,
) dto.RegistrationResponse {

	result := dto.RegistrationResponse{
		ID:            r.ID,
		WorkshopID:    r.WorkshopID,
		AttendeeName:  r.AttendeeName,
		AttendeeEmail: r.AttendeeEmail,
		Status:        r.Status,
		RegisteredBy:  r.RegisteredBy,
		RegisteredAt:  r.RegisteredAt.Format(time.RFC3339),
		CancelledBy:   r.CancelledBy,
	}

	if r.CancelledAt != nil {
		formatted := r.CancelledAt.Format(time.RFC3339)
		result.CancelledAt = &formatted
	}

	return result
}

func handleRegistrationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(404, gin.H{"error": "Record not found"})

	case errors.Is(err, service.ErrWorkshopFull),
		errors.Is(err, service.ErrWorkshopClosed),
		errors.Is(err, service.ErrAlreadyCancelled):
		c.JSON(409, gin.H{"error": err.Error()})

	case errors.Is(err, service.ErrInvalidAttendee):
		c.JSON(400, gin.H{"error": err.Error()})

	default:
		c.JSON(500, gin.H{"error": "Internal server error"})
	}
}

func (h *RegistrationHandler) Create(c *gin.Context) {
	id, ok := workshopID(c)
	if !ok {
		return
	}

	var req dto.RegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	staffID := c.GetUint("userID")

	registration, err := h.Service.Create(
		id,
		req.AttendeeName,
		req.AttendeeEmail,
		staffID,
	)

	if err != nil {
		handleRegistrationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, registrationResponse(*registration))
}
func (h *RegistrationHandler) GetByWorkshop(c *gin.Context) {
	id, ok := workshopID(c)
	if !ok {
		return
	}

	registrations, err := h.Service.GetByWorkshop(id)
	if err != nil {
		handleRegistrationError(c, err)
		return
	}

	result := make([]dto.RegistrationResponse, 0, len(registrations))

	for _, registration := range registrations {
		result = append(result, registrationResponse(registration))
	}

	c.JSON(http.StatusOK, result)
}

func (h *RegistrationHandler) Cancel(c *gin.Context) {
	id, ok := workshopID(c)
	if !ok {
		return
	}

	staffID := c.GetUint("userID")

	registration, err := h.Service.Cancel(id, staffID)

	if err != nil {
		handleRegistrationError(c, err)
		return
	}

	c.JSON(http.StatusOK, registrationResponse(*registration))
}
