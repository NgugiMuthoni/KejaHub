package routes

import (
	"github.com/gin-gonic/gin"

	"kejahub-backend/internal/database"
)

type CreateHouseInput struct {
	Name       string `json:"name"`
	Location   string `json:"location"`
	TotalUnits int    `json:"total_units"`
}

func CreateHouse(c *gin.Context) {

	role := c.GetString("role")

	if role != "landlord" && role != "admin" {
		c.JSON(403, gin.H{"error": "only landlords can create houses"})
		return
	}

	email := c.GetString("email")

	users, err := database.Get("users", "?email=eq."+email)
	if err != nil || len(users) == 0 {
		c.JSON(500, gin.H{"error": "user not found"})
		return
	}

	userID := users[0]["id"].(string)

	landlords, err := database.Get("landlords", "?user_id=eq."+userID)
	if err != nil || len(landlords) == 0 {
		c.JSON(500, gin.H{"error": "landlord not found"})
		return
	}

	landlordID := landlords[0]["id"].(string)

	var input CreateHouseInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"landlord_id": landlordID,
		"name":        input.Name,
		"location":    input.Location,
		"total_units": input.TotalUnits,
	}

	err = database.Insert("houses", data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "house created successfully"})
}
