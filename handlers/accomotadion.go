package handlers

import (
	"net/http"
	"go-backend-basic/config"

	"github.com/gin-gonic/gin"
)

//
// =======================
// CREATE UNIT (accommodations + unit_details + galleries)
// =======================
func CreateAccommodation(c *gin.Context) {
	var input struct {
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Price       float64 `json:"price"`
		Alamat      string  `json:"alamat"`
		JumlahKamar int     `json:"jumlah_kamar"`
		Fasilitas   string  `json:"fasilitas"`
		ImageURL    string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 1️⃣ insert accommodations
	var accommodationID int
	err = tx.QueryRow(`
		INSERT INTO accommodations (name, type, price)
		VALUES ($1,$2,$3)
		RETURNING id
	`, input.Name, input.Type, input.Price).Scan(&accommodationID)

	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ insert unit_details
	_, err = tx.Exec(`
		INSERT INTO unit_details (accommodation_id, alamat, jumlah_kamar, fasilitas)
		VALUES ($1,$2,$3,$4)
	`, accommodationID, input.Alamat, input.JumlahKamar, input.Fasilitas)

	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 3️⃣ insert gallery (optional)
	if input.ImageURL != "" {
		_, err = tx.Exec(`
			INSERT INTO galleries (accommodation_id, image_url)
			VALUES ($1,$2)
		`, accommodationID, input.ImageURL)

		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Unit berhasil ditambahkan",
	})
}

//
// =======================
// READ ALL UNITS
// =======================
func GetAccommodations(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT a.id, a.name, a.type, a.price,
		       COALESCE(d.alamat, '') AS alamat,
		       COALESCE(d.jumlah_kamar, 0) AS jumlah_kamar,
		       COALESCE(d.fasilitas, '') AS fasilitas
		FROM accommodations a
		LEFT JOIN unit_details d ON d.accommodation_id = a.id
		ORDER BY a.id DESC
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	data := make([]gin.H, 0)

	for rows.Next() {
		var (
			id           int
			name, typ    string
			alamat       string
			fasilitas    string
			jumlahKamar  int
			price        float64
		)

		err := rows.Scan(
			&id, &name, &typ, &price,
			&alamat, &jumlahKamar, &fasilitas,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		data = append(data, gin.H{
			"id":            id,
			"name":          name,
			"type":          typ,
			"price":         price,
			"alamat":        alamat,
			"jumlah_kamar":  jumlahKamar,
			"fasilitas":     fasilitas,
		})
	}

	c.JSON(http.StatusOK, data)
}

//
// =======================
// READ UNIT BY ID
// =======================
func GetAccommodationByID(c *gin.Context) {
	id := c.Param("id")

	var (
		accID        int
		name         string
		typ          string
		price        float64
		alamat       string
		jumlahKamar  int
		fasilitas    string
	)

	err := config.DB.QueryRow(`
		SELECT a.id, a.name, a.type, a.price,
		       COALESCE(d.alamat, ''),
		       COALESCE(d.jumlah_kamar, 0),
		       COALESCE(d.fasilitas, '')
		FROM accommodations a
		LEFT JOIN unit_details d ON d.accommodation_id = a.id
		WHERE a.id = $1
	`, id).Scan(
		&accID, &name, &typ, &price,
		&alamat, &jumlahKamar, &fasilitas,
	)

	if err != nil {
		c.JSON(404, gin.H{"message": "Unit tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            accID,
		"name":          name,
		"type":          typ,
		"price":         price,
		"alamat":        alamat,
		"jumlah_kamar":  jumlahKamar,
		"fasilitas":     fasilitas,
	})
}

//
// =======================
// UPDATE UNIT
// =======================
func UpdateAccommodation(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Price       float64 `json:"price"`
		Alamat      string  `json:"alamat"`
		JumlahKamar int     `json:"jumlah_kamar"`
		Fasilitas   string  `json:"fasilitas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Exec(`
		UPDATE accommodations
		SET name=$1, type=$2, price=$3
		WHERE id=$4
	`, input.Name, input.Type, input.Price, id)

	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Exec(`
		UPDATE unit_details
		SET alamat=$1, jumlah_kamar=$2, fasilitas=$3
		WHERE accommodation_id=$4
	`, input.Alamat, input.JumlahKamar, input.Fasilitas, id)

	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "Unit berhasil diupdate"})
}

//
// =======================
// DELETE UNIT
// =======================
func DeleteAccommodation(c *gin.Context) {
	id := c.Param("id")

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Exec(`DELETE FROM galleries WHERE accommodation_id=$1`, id)
	tx.Exec(`DELETE FROM unit_details WHERE accommodation_id=$1`, id)

	_, err = tx.Exec(`DELETE FROM accommodations WHERE id=$1`, id)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "Unit berhasil dihapus"})
}
