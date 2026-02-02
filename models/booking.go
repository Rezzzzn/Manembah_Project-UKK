package models

import "time"

type Booking struct {
	IDBooking     int       `json:"id_booking" db:"id_booking"`
	UserID        int       `json:"user_id" db:"user_id"`
	UnitID        int       `json:"unit_id" db:"unit_id"`
	CheckIn       time.Time `json:"check_in" db:"check_in"`
	CheckOut      time.Time `json:"check_out" db:"check_out"`
	JumlahOrang   int       `json:"jumlah_orang" db:"jumlah_orang"`
	StatusBooking string    `json:"status_booking" db:"status_booking"`
	InvoiceNumber string    `json:"invoice_number" db:"invoice_number"` // ⬅️ BARU
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
