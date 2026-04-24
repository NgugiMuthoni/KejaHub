package routes

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"kejahub-backend/internal/database"
)

type CreateTenantInput struct {
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	UnitNumber string  `json:"unit_number"`
	RentAmount float64 `json:"rent_amount"`
	HouseID    string  `json:"house_id"` // still needed BUT validated via JWT
}

func CreateTenant(c *gin.Context) {

	role := c.GetString("role")
	email := c.GetString("email")

	fmt.Println("ROLE FROM CONTEXT:", c.GetString("role"))
	fmt.Println("EMAIL FROM CONTEXT:", c.GetString("email"))

	// 1. Only landlords/admins can create tenants
	if role != "landlord" && role != "admin" {
		c.JSON(403, gin.H{"error": "only landlords can create tenants"})
		return
	}

	// 2. Get user
	users, err := database.Get("users", "?email=eq."+email)
	if err != nil || len(users) == 0 {
		c.JSON(500, gin.H{"error": "user not found"})
		return
	}

	userID := users[0]["id"].(string)

	// 3. Get landlord
	landlords, err := database.Get("landlords", "?user_id=eq."+userID)
	if err != nil || len(landlords) == 0 {
		c.JSON(500, gin.H{"error": "landlord not found"})
		return
	}

	landlordID := landlords[0]["id"].(string)

	// 4. Verify house belongs to this landlord (IMPORTANT SECURITY STEP)
	var input CreateTenantInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	houses, err := database.Get("houses", "?id=eq."+input.HouseID+"&landlord_id=eq."+landlordID)
	if err != nil || len(houses) == 0 {
		c.JSON(403, gin.H{"error": "invalid house or unauthorized access"})
		return
	}

	// 5. Create tenant (SAFE)
	data := map[string]interface{}{
		"house_id":    input.HouseID,
		"name":        input.Name,
		"phone":       input.Phone,
		"unit_number": input.UnitNumber,
		"rent_amount": input.RentAmount,
	}

	err = database.Insert("tenants", data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "tenant created successfully",
	})
}
