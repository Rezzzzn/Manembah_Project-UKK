package services

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	from := "yourmail@gmail.com"     // ganti
	password := "APP_PASSWORD_GMAIL" // app password

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	return smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
}

func BuildInvoiceEmail(invoice, unit, checkIn, checkOut string) string {
	return fmt.Sprintf(`
<h2>Invoice Booking</h2>
<p>Invoice: %s</p>
<p>Unit: %s</p>
<p>Check-in: %s</p>
<p>Check-out: %s</p>
<p>Status: UNPAID</p>
<p>Silakan selesaikan pembayaran sebelum H-7.</p>
`, invoice, unit, checkIn, checkOut)
}

func BuildReminderEmail(invoice, checkIn string) string {
	return fmt.Sprintf(`
<h2>Reminder Pembayaran Booking</h2>
<p>Invoice: %s</p>
<p>Check-in: %s</p>
<p>Booking Anda belum dibayar.</p>
<p>Jika tidak dibayar sebelum H-7, booking akan dibatalkan otomatis.</p>
`, invoice, checkIn)
}

func BuildRefundEmail(invoice string) string {
	return fmt.Sprintf(`
<h2>Refund Booking</h2>
<p>Invoice: %s</p>
<p>Booking Anda telah dibatalkan.</p>
<p>Dana akan dikembalikan sesuai kebijakan refund.</p>
`, invoice)
}
