package dto

type RegistrationRequest struct {
	AttendeeName  string `json:"attendeeName" binding:"required"`
	AttendeeEmail string `json:"attendeeEmail" binding:"required,email"`
}

type RegistrationResponse struct {
	ID            uint    `json:"id"`
	WorkshopID    uint    `json:"workshopId"`
	AttendeeName  string  `json:"attendeeName"`
	AttendeeEmail string  `json:"attendeeEmail"`
	Status        string  `json:"status"`
	RegisteredBy  uint    `json:"registeredBy"`
	RegisteredAt  string  `json:"registeredAt"`
	CancelledBy   *uint   `json:"cancelledBy"`
	CancelledAt   *string `json:"cancelledAt"`
}
