package dto

type PaymentResponse struct {
	PaymentID   int    `json:"payment_id"`
	BookingID   int    `json:"booking_id"`
	PaymentURL  string `json:"payment_url"`
	PaymentType string `json:"payment_type"`
	Status      string `json:"status"`
}
