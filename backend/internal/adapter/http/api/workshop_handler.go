package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"workshop-registration/backend/internal/adapter/http/dto"
	"workshop-registration/backend/internal/adapter/storage/postgres/models"
	"workshop-registration/backend/internal/core/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkshopHandler struct {
	Service *service.WorkshopService
}

func NewWorkshopHandler(
	s *service.WorkshopService,
) *WorkshopHandler {
	return &WorkshopHandler{Service: s}
}

func workshopID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 || id > uint64(^uint(0)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workshop ID",
		})
		return 0, false
	}

	return uint(id), true
}

func workshopError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Workshop not found",
		})
	case errors.Is(err, gorm.ErrDuplicatedKey):
		c.JSON(http.StatusConflict, gin.H{
			"error": "Workshop code already exists",
		})
	case errors.Is(err, service.ErrCapacityTooLow),
		errors.Is(err, service.ErrWorkshopInvalid):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Request failed",
		})
	}
}

func (h *WorkshopHandler) response(
	w models.WorkshopModel,
) (dto.WorkshopResponse, error) {

	count, err := h.Service.Repo.ActiveCount(w.ID)
	if err != nil {
		return dto.WorkshopResponse{}, err
	}

	return dto.WorkshopResponse{
		ID:             w.ID,
		Code:           w.Code,
		Title:          w.Title,
		Instructor:     w.Instructor,
		StartsAt:       w.StartsAt.Format(time.RFC3339),
		Capacity:       w.Capacity,
		Status:         w.Status,
		AvailableSeats: w.Capacity - int(count),
	}, nil
}

func (h *WorkshopHandler) Create(c *gin.Context) {
	var req dto.WorkshopRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workshop, err := h.Service.Create(req)
	if err != nil {
		workshopError(c, err)
		return
	}

	result, err := h.response(*workshop)
	if err != nil {
		workshopError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *WorkshopHandler) GetAll(c *gin.Context) {

	from := c.Query("from")
	to := c.Query("to")
	status := c.Query("status")

	availableOnly := false

	if value, exists := c.GetQuery("availableOnly"); exists {
		parsed, err := strconv.ParseBool(value)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "availableOnly must be true or false",
			})
			return
		}

		availableOnly = parsed
	}

	workshops, err := h.Service.GetFiltered(
		from,
		to,
		status,
		availableOnly,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidFilter) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date range or status filter",
			})
			return
		}

		workshopError(c, err)
		return
	}

	result := make(
		[]dto.WorkshopResponse,
		0,
		len(workshops),
	)

	for _, workshop := range workshops {

		item, err := h.response(workshop)

		if err != nil {
			workshopError(c, err)
			return
		}

		result = append(result, item)
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkshopHandler) GetByID(c *gin.Context) {
	id, ok := workshopID(c)
	if !ok {
		return
	}

	workshop, err := h.Service.GetByID(id)
	if err != nil {
		workshopError(c, err)
		return
	}

	result, err := h.response(*workshop)
	if err != nil {
		workshopError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkshopHandler) Update(c *gin.Context) {
	id, ok := workshopID(c)
	if !ok {
		return
	}

	var req dto.WorkshopRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workshop, err := h.Service.Update(id, req)
	if err != nil {
		workshopError(c, err)
		return
	}

	result, err := h.response(*workshop)
	if err != nil {
		workshopError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
