package handlers

import (
	"net/http"
	"go-backend-basic/config"
	"github.com/gin-gonic/gin"
)

func CreateAccommodation(c *gin.Context) {
	var input struct {
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		Location    string  `json:"location"`
		Price       float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO accommodations (name, type, description, location, price)
		VALUES ($1,$2,$3,$4,$5)
	`, input.Name, input.Type, input.Description, input.Location, input.Price)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Data berhasil ditambahkan"})
}

// read

func GetAccommodations(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, type, description, location, price FROM accommodations")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var data []gin.H

	for rows.Next() {
		var id int
		var name, typ, desc, loc string
		var price float64

		rows.Scan(&id, &name, &typ, &desc, &loc, &price)

		data = append(data, gin.H{
			"id": id,
			"name": name,
			"type": typ,
			"description": desc,
			"location": loc,
			"price": price,
		})
	}

	c.JSON(200, data)
}


// update
func UpdateAccommodation(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Location    string  `json:"location"`
		Price       float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(`
		UPDATE accommodations
		SET name=$1, description=$2, location=$3, price=$4, updated_at=NOW()
		WHERE id=$5
	`, input.Name, input.Description, input.Location, input.Price, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Data berhasil diupdate"})
}


// delete
func DeleteAccommodation(c *gin.Context) {
	id := c.Param("id")

	_, err := config.DB.Exec("DELETE FROM accommodations WHERE id=$1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Data berhasil dihapus"})
}


