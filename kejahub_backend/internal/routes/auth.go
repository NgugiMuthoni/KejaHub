package routes

import (
	"github.com/gin-gonic/gin"

	"kejahub-backend/internal/auth"
	"kejahub-backend/internal/database"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// REGISTER (auto landlord + role)
func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := auth.HashPassword(input.Password)

	// 1. create user with role
	userData := map[string]interface{}{
		"email":    input.Email,
		"password": hashedPassword,
		"role":     "landlord",
		"name":     input.Name,
		"phone":    input.Phone,
	}

	err := database.Insert("users", userData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 2. fetch user
	users, err := database.Get("users", "?email=eq."+input.Email)
	if err != nil || len(users) == 0 {
		c.JSON(500, gin.H{"error": "user fetch failed"})
		return
	}

	userID := users[0]["id"].(string)

	// 3. auto create landlord
	landlordData := map[string]any{
		"user_id": userID,
		"name":    input.Name,
		"phone":   input.Phone,
	}

	err = database.Insert("landlords", landlordData)
	if err != nil {
		c.JSON(500, gin.H{"error": "landlord creation failed"})
		return
	}

	c.JSON(200, gin.H{
		"message": "user + landlord created successfully",
	})
}

// LOGIN (returns JWT with role)
func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	users, err := database.Get("users", "?email=eq."+input.Email)
	if err != nil || len(users) == 0 {
		c.JSON(401, gin.H{"error": "user not found"})
		return
	}

	user := users[0]

	if !auth.CheckPassword(input.Password, user["password"].(string)) {
		c.JSON(401, gin.H{"error": "invalid password"})
		return
	}

	email := user["email"].(string)
	role := user["role"].(string)

	token, _ := auth.GenerateToken(email, role)

	c.JSON(200, gin.H{
		"token": token,
	})
}
