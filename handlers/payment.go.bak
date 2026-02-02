package handlers

import (
	"net/http"
	"time"

	"go-backend-basic/config"

	"github.com/gin-gonic/gin"
)

// ========================
// CREATE PAYMENT (DP)
// ========================
func CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek booking exist & status
	var status string
	err := config.DB.QueryRow(`
		SELECT status_booking 
		FROM booking 
		WHERE id_booking = $1
	`, req.BookingID).Scan(&status)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	if status == "cancelled" || status == "refunded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking not valid for payment"})
		return
	}

	// insert payment
	query := `
		INSERT INTO payment (booking_id, amount, status_payment)
		VALUES ($1,$2,'pending')
		RETURNING payment_id
	`

	var paymentID int
	err = config.DB.QueryRow(query, req.BookingID, req.Amount).Scan(&paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ========================
	// MIDTRANS INTEGRATION (STUB)
	// ========================
	paymentURL := "https://sandbox.midtrans.com/snap/v2/vtweb/dummy-payment-url"

	// update booking status
	_, _ = config.DB.Exec(`
		UPDATE booking 
		SET status_booking = 'dp_pending' 
		WHERE id_booking = $1
	`, req.BookingID)

	c.JSON(http.StatusCreated, PaymentResponse{
		PaymentID:   paymentID,
		BookingID:   req.BookingID,
		PaymentURL:  paymentURL,
		PaymentType: "dp",
		Status:      "pending",
	})
}

// ========================
// GET PAYMENT BY BOOKING
// ========================
func GetPaymentByBooking(c *gin.Context) {
	bookingID := c.Param("booking_id")

	row := config.DB.QueryRow(`
		SELECT payment_id, booking_id, amount, status_payment
		FROM payment
		WHERE booking_id = $1
		ORDER BY payment_id DESC
		LIMIT 1
	`, bookingID)

	var paymentID int
	var bID int
	var amount float64
	var status string

	err := row.Scan(&paymentID, &bID, &amount, &status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_id": paymentID,
		"booking_id": bID,
		"amount":     amount,
		"status":     status,
	})
}

// ========================
// UPDATE PAYMENT STATUS (WEBHOOK SIMULATION)
// ========================
func UpdatePaymentStatus(c *gin.Context) {
	var body struct {
		PaymentID int    `json:"payment_id"`
		Status    string `json:"status_payment"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(`
		UPDATE payment 
		SET status_payment = $1 
		WHERE payment_id = $2
	`, body.Status, body.PaymentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// kalau settlement → update booking
	if body.Status == "settlement" {
		_, _ = config.DB.Exec(`
			UPDATE booking 
			SET status_booking = 'dp_paid'
			WHERE id_booking = (
				SELECT booking_id FROM payment WHERE payment_id = $1
			)
		`, body.PaymentID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment status updated"})
}

// ========================
// REFUND PAYMENT (H-7 RULE)
// ========================
func RefundPayment(c *gin.Context) {
	var body struct {
		BookingID int `json:"booking_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get booking
	var checkIn time.Time
	err := config.DB.QueryRow(`
		SELECT check_in 
		FROM booking 
		WHERE id_booking = $1
	`, body.BookingID).Scan(&checkIn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	diff := checkIn.Sub(time.Now()).Hours() / 24
	if diff < 7 {
		c.JSON(http.StatusForbidden, gin.H{"error": "refund only allowed before H-7"})
		return
	}

	// update payment
	_, err = config.DB.Exec(`
		UPDATE payment 
		SET status_payment = 'refund' 
		WHERE booking_id = $1
	`, body.BookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update booking
	_, _ = config.DB.Exec(`
		UPDATE booking 
		SET status_booking = 'refunded' 
		WHERE id_booking = $1
	`, body.BookingID)

	c.JSON(http.StatusOK, gin.H{
		"message": "refund processed",
		"status":  "refunded",
	})
}
