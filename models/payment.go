package models

import "time"

type Payment struct {
	PaymentID int       `json:"payment_id" db:"payment_id"`
	BookingID int       `json:"booking_id" db:"booking_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Status    string    `json:"status_payment" db:"status_payment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
