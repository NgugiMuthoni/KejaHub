package routes

import (
	"time"

	"kejahub-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func GenerateMonthlyRent(c *gin.Context) {

	tenants, _ := database.Get("tenants", "")

	currentMonth := time.Now().Format("2006-01")

	for _, t := range tenants {

		tenantID := t["id"].(string)
		rent := t["rent_amount"].(float64)

		// check if already exists
		existing, _ := database.Get("rent_cycles", "?tenant_id=eq."+tenantID+"&month=eq."+currentMonth)

		if len(existing) > 0 {
			continue
		}

		data := map[string]interface{}{
			"tenant_id":   tenantID,
			"month":       currentMonth,
			"rent_amount": rent,
			"due_date":    time.Now(),
		}

		database.Insert("rent_cycles", data)
	}

	c.JSON(200, gin.H{"message": "rent generated"})
}
