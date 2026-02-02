package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go-backend-basic/config"
	"go-backend-basic/services"

	"github.com/gin-gonic/gin"
)

// ========================
// STRUCTS
// ========================

type CreateBookingRequest struct {
	UnitID      int    `json:"unit_id" binding:"required"`
	CheckIn     string `json:"check_in" binding:"required"`
	CheckOut    string `json:"check_out" binding:"required"`
	JumlahOrang int    `json:"jumlah_orang"`
}

// ========================
// HELPERS
// ========================

func GenerateInvoiceNumber() string {
	return "INV-" + time.Now().Format("20060102-150405")
}

// ========================
// CREATE BOOKING
// ========================

func CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// ambil email user
	var userEmail string
	err := config.DB.QueryRow(
		`SELECT email FROM users WHERE id_user = $1`,
		userID,
	).Scan(&userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed get user email"})
		return
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_in format (YYYY-MM-DD)"})
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_out format (YYYY-MM-DD)"})
		return
	}

	invoice := GenerateInvoiceNumber()

	// === INSERT BOOKING ===
	var bookingID int
	err = config.DB.QueryRow(`
		INSERT INTO booking 
		(user_id, unit_id, check_in, check_out, jumlah_orang, status_booking, invoice_number)
		VALUES ($1,$2,$3,$4,$5,'pending',$6)
		RETURNING id_booking
	`,
		userID,
		req.UnitID,
		checkIn,
		checkOut,
		req.JumlahOrang,
		invoice,
	).Scan(&bookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// === HITUNG TOTAL HARGA ===
	var hargaPerMalam int
	err = config.DB.QueryRow(
		`SELECT harga_per_malam FROM unit WHERE id_unit = $1`,
		req.UnitID,
	).Scan(&hargaPerMalam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed get unit price"})
		return
	}

	durasi := int(checkOut.Sub(checkIn).Hours() / 24)
	if durasi <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
		return
	}

	totalHarga := hargaPerMalam * durasi

	// === INSERT PAYMENT ===
	_, err = config.DB.Exec(`
		INSERT INTO payment 
		(booking_id, invoice_number, amount, payment_status)
		VALUES ($1,$2,$3,'unpaid')
	`, bookingID, invoice, totalHarga)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// =========================
	// 📧 SEND EMAIL INVOICE
	// =========================
	body := services.BuildInvoiceEmail(
		invoice,
		"Unit Booking",
		req.CheckIn,
		req.CheckOut,
	)

	go services.SendEmail(userEmail, "Invoice Booking", body) // async

	c.JSON(http.StatusCreated, gin.H{
		"message":    "booking created",
		"id_booking": bookingID,
		"invoice":    invoice,
		"status":     "pending",
		"payment":    "unpaid",
	})
}

// ========================
// GET USER BOOKINGS
// ========================

func GetMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := config.DB.Query(`
		SELECT 
			b.id_booking, b.unit_id, b.check_in, b.check_out, 
			b.jumlah_orang, b.status_booking, b.invoice_number,
			p.payment_status, p.amount
		FROM booking b
		JOIN payment p ON p.booking_id = b.id_booking
		WHERE b.user_id = $1
		ORDER BY b.id_booking DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type BookingResponse struct {
		IDBooking   int       `json:"id_booking"`
		UnitID      int       `json:"unit_id"`
		CheckIn     time.Time `json:"check_in"`
		CheckOut    time.Time `json:"check_out"`
		JumlahOrang int       `json:"jumlah_orang"`
		Status      string    `json:"status_booking"`
		Invoice     string    `json:"invoice"`
		Payment     string    `json:"payment_status"`
		Amount      int       `json:"amount"`
	}

	bookings := []BookingResponse{}

	for rows.Next() {
		var b BookingResponse
		err := rows.Scan(
			&b.IDBooking,
			&b.UnitID,
			&b.CheckIn,
			&b.CheckOut,
			&b.JumlahOrang,
			&b.Status,
			&b.Invoice,
			&b.Payment,
			&b.Amount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookings = append(bookings, b)
	}

	c.JSON(http.StatusOK, bookings)
}

// ========================
// GET ALL BOOKINGS (ADMIN)
// ========================

func GetAllBookings(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT 
			b.id_booking, b.user_id, b.unit_id, b.check_in, b.check_out,
			b.jumlah_orang, b.status_booking, b.invoice_number,
			p.payment_status, p.amount
		FROM booking b
		JOIN payment p ON p.booking_id = b.id_booking
		ORDER BY b.id_booking DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	bookings := []gin.H{}

	for rows.Next() {
		var id, userID, unitID, jumlahOrang, amount int
		var checkIn, checkOut time.Time
		var status, invoice, payment string

		rows.Scan(
			&id, &userID, &unitID, &checkIn, &checkOut,
			&jumlahOrang, &status, &invoice, &payment, &amount,
		)

		bookings = append(bookings, gin.H{
			"id_booking":     id,
			"user_id":        userID,
			"unit_id":        unitID,
			"check_in":       checkIn,
			"check_out":      checkOut,
			"jumlah_orang":   jumlahOrang,
			"status_booking": status,
			"invoice":        invoice,
			"payment_status": payment,
			"amount":         amount,
		})
	}

	c.JSON(http.StatusOK, bookings)
}

// ========================
// UPDATE BOOKING STATUS (ADMIN)
// ========================

func UpdateBookingStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status_booking"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingID, _ := strconv.Atoi(id)

	_, err := config.DB.Exec(`
		UPDATE booking
		SET status_booking = $1
		WHERE id_booking = $2
	`, body.Status, bookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking status updated"})
}

// ========================
// CANCEL BOOKING (H-7 RULE)
// ========================

func CancelBooking(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	bookingID, _ := strconv.Atoi(id)

	var checkIn time.Time
	var status string
	var paymentStatus string

	err := config.DB.QueryRow(`
		SELECT b.check_in, b.status_booking, p.payment_status
		FROM booking b
		JOIN payment p ON p.booking_id = b.id_booking
		WHERE b.id_booking = $1 AND b.user_id = $2
	`, bookingID, userID).Scan(&checkIn, &status, &paymentStatus)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	if status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking already cancelled"})
		return
	}

	now := time.Now()
	diff := checkIn.Sub(now).Hours() / 24

	// RULE H-7
	if diff < 7 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "refund only allowed before H-7 check-in",
		})
		return
	}

	// update booking
	_, err = config.DB.Exec(`
		UPDATE booking
		SET status_booking = 'cancelled'
		WHERE id_booking = $1
	`, bookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update payment
	_, err = config.DB.Exec(`
		UPDATE payment
		SET payment_status = 'refunded'
		WHERE booking_id = $1
	`, bookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking cancelled",
		"refund":  "eligible",
		"payment": "refunded",
	})
}
