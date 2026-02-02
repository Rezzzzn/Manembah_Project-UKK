package cron

import (
	"time"

	"go-backend-basic/config"
	"go-backend-basic/services"
)

func StartScheduler() {
	go H7ReminderJob()
	go AutoCancelUnpaidJob()
}

// H-7 Reminder
func H7ReminderJob() {
	for {
		rows, _ := config.DB.Query(`
SELECT u.email, b.invoice_number, b.check_in
FROM booking b
JOIN payment p ON p.booking_id = b.id_booking
JOIN users u ON u.id_user = b.user_id
WHERE p.payment_status = 'unpaid'
AND b.status_booking = 'pending'
AND b.check_in - NOW() <= INTERVAL '7 days'
`)

		for rows.Next() {
			var email, invoice string
			var checkIn time.Time

			rows.Scan(&email, &invoice, &checkIn)

			body := services.BuildReminderEmail(invoice, checkIn.Format("2006-01-02"))
			services.SendEmail(email, "Reminder Pembayaran Booking", body)
		}
		rows.Close()

		time.Sleep(24 * time.Hour)
	}
}

// Auto Cancel Unpaid
func AutoCancelUnpaidJob() {
	for {
		rows, _ := config.DB.Query(`
SELECT b.id_booking, u.email, b.invoice_number
FROM booking b
JOIN payment p ON p.booking_id = b.id_booking
JOIN users u ON u.id_user = b.user_id
WHERE p.payment_status = 'unpaid'
AND b.check_in <= NOW()
AND b.status_booking = 'pending'
`)

		for rows.Next() {
			var id int
			var email, invoice string

			rows.Scan(&id, &email, &invoice)

			config.DB.Exec(`UPDATE booking SET status_booking='cancelled' WHERE id_booking=$1`, id)
			config.DB.Exec(`UPDATE payment SET payment_status='expired' WHERE booking_id=$1`, id)

			body := services.BuildRefundEmail(invoice)
			services.SendEmail(email, "Booking Dibatalkan", body)
		}
		rows.Close()

		time.Sleep(24 * time.Hour)
	}
}
