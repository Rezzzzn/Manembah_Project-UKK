package dto

type CreatePaymentRequest struct {
	BookingID int     `json:"booking_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}
