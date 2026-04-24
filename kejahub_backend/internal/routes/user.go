package routes

import (
	"kejahub-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {

	email := c.GetString("email")

	users, err := database.Get("users", "?email=eq."+email)
	if err != nil || len(users) == 0 {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, users[0])
}
