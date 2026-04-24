package routes

import (
	"kejahub-backend/internal/database"

	"github.com/gin-gonic/gin"
)

type AddUtilityInput struct {
	RentCycleID string `json:"rent_cycle_id"`
	Name        string `json:"name"`
	IsPriority  bool   `json:"is_priority"`

	Amount          float64 `json:"amount"`
	PreviousReading float64 `json:"previous_reading"`
	CurrentReading  float64 `json:"current_reading"`
	UnitCost        float64 `json:"unit_cost"`
}

func AddUtility(c *gin.Context) {

	role := c.GetString("role")

	if role != "landlord" && role != "caretaker" && role != "admin" {
		c.JSON(403, gin.H{"error": "not allowed"})
		return
	}

	var input AddUtilityInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var finalAmount float64 = input.Amount
	var unitsUsed float64 = 0

	// 🧠 Meter-based calculation
	if input.CurrentReading > 0 && input.UnitCost > 0 {

		unitsUsed = input.CurrentReading - input.PreviousReading

		if unitsUsed < 0 {
			c.JSON(400, gin.H{"error": "invalid readings"})
			return
		}

		finalAmount = unitsUsed * input.UnitCost
	}

	data := map[string]interface{}{
		"rent_cycle_id":    input.RentCycleID,
		"name":             input.Name,
		"amount":           finalAmount,
		"is_priority":      input.IsPriority,
		"previous_reading": input.PreviousReading,
		"current_reading":  input.CurrentReading,
		"unit_cost":        input.UnitCost,
		"units_used":       unitsUsed,
	}

	err := database.Insert("utilities", data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "utility added"})
}
