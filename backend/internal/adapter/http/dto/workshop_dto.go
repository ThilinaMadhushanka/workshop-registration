package dto

type WorkshopRequest struct {
	Code       string `json:"code" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Instructor string `json:"instructor" binding:"required"`
	StartsAt   string `json:"startsAt" binding:"required"`
	Capacity   int    `json:"capacity" binding:"required,gt=0"`
	Status     string `json:"status" binding:"required,oneof=scheduled cancelled completed"`
}

type WorkshopResponse struct {
	ID             uint   `json:"id"`
	Code           string `json:"code"`
	Title          string `json:"title"`
	Instructor     string `json:"instructor"`
	StartsAt       string `json:"startsAt"`
	Capacity       int    `json:"capacity"`
	Status         string `json:"status"`
	AvailableSeats int    `json:"availableSeats"`
}
