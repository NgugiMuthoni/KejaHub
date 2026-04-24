package routes

import (
	"kejahub-backend/internal/auth"
	"kejahub-backend/internal/database"

	"math/rand"

	"github.com/gin-gonic/gin"
)

type CreateCaretakerInput struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	HouseID string `json:"house_id"`
}

func generatePassword() string {
	return "pass" + string(rand.Intn(100000))
}

func CreateCaretaker(c *gin.Context) {

	role := c.GetString("role")
	email := c.GetString("email")
	password := generatePassword()

	hashed, _ := auth.HashPassword(password)

	if role != "landlord" && role != "admin" {
		c.JSON(403, gin.H{"error": "only landlords can create caretakers"})
		return
	}

	var input CreateCaretakerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 1. get landlord
	users, _ := database.Get("users", "?email=eq."+email)
	userID := users[0]["id"].(string)

	landlords, _ := database.Get("landlords", "?user_id=eq."+userID)
	landlordID := landlords[0]["id"].(string)

	// 2. verify house belongs to landlord
	houses, _ := database.Get("houses", "?id=eq."+input.HouseID+"&landlord_id=eq."+landlordID)
	if len(houses) == 0 {
		c.JSON(403, gin.H{"error": "invalid house"})
		return
	}

	// 3. create caretaker user
	userData := map[string]interface{}{
		"email":    input.Email,
		"password": hashed,
		"role":     "caretaker",
		"name":     input.Name,
		"phone":    input.Phone,
	}

	database.Insert("users", userData)

	// 4. fetch new user
	newUser, _ := database.Get("users", "?email=eq."+input.Email)
	newUserID := newUser[0]["id"].(string)

	// 5. insert caretaker
	data := map[string]interface{}{
		"user_id":  newUserID,
		"house_id": input.HouseID,
		"name":     input.Name,
		"phone":    input.Phone,
	}

	database.Insert("caretakers", data)

	c.JSON(200, gin.H{"message": "caretaker created"})
}
