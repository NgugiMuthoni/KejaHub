package routes

import (
	"github.com/gin-gonic/gin"
	"kejahub-backend/internal/database"
)

func CheckPriorityUtilities(c *gin.Context) {

	rentCycleID := c.Param("id")

	utilities, _ := database.Get(
		"utilities",
		"?rent_cycle_id=eq."+rentCycleID+"&is_priority=eq.true",
	)

	if len(utilities) == 0 {
		c.JSON(400, gin.H{
			"error": "priority utilities not completed",
		})
		return
	}

	c.JSON(200, gin.H{"message": "ready"})
}