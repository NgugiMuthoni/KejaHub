package routes

import (
	"kejahub-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func LandlordDashboard(c *gin.Context) {

	role := c.GetString("role")
	email := c.GetString("email")

	if role != "landlord" && role != "admin" {
		c.JSON(403, gin.H{"error": "not allowed"})
		return
	}

	// 1. get user
	users, _ := database.Get("users", "?email=eq."+email)
	userID := users[0]["id"].(string)

	// 2. get landlord
	landlords, _ := database.Get("landlords", "?user_id=eq."+userID)
	landlordID := landlords[0]["id"].(string)

	// 3. get houses
	houses, _ := database.Get("houses", "?landlord_id=eq."+landlordID)

	totalHouses := len(houses)

	totalTenants := 0
	totalRevenue := 0.0
	unpaidTenants := 0

	// 4. loop houses → tenants → payments
	for _, house := range houses {

		houseID := house["id"].(string)

		tenants, _ := database.Get("tenants", "?house_id=eq."+houseID)

		totalTenants += len(tenants)

		for _, tenant := range tenants {

			tenantID := tenant["id"].(string)
			rent := tenant["rent_amount"].(float64)

			payments, _ := database.Get("payments", "?tenant_id=eq."+tenantID)

			var paid float64 = 0

			for _, p := range payments {
				paid += p["amount"].(float64)
			}

			totalRevenue += paid

			if paid < rent {
				unpaidTenants++
			}
		}
	}

	c.JSON(200, gin.H{
		"total_houses":   totalHouses,
		"total_tenants":  totalTenants,
		"total_revenue":  totalRevenue,
		"unpaid_tenants": unpaidTenants,
	})
}
