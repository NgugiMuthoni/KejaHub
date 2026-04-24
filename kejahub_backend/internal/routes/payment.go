package routes

import (
	"time"

	"kejahub-backend/internal/database"

	"github.com/gin-gonic/gin"
)

type CreatePaymentInput struct {
	RentCycleID string  `json:"rent_cycle_id"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
}

func CreatePayment(c *gin.Context) {

	role := c.GetString("role")
	email := c.GetString("email")

	// 🔐 Only allowed roles
	if role != "landlord" && role != "caretaker" && role != "admin" {
		c.JSON(403, gin.H{"error": "not allowed"})
		return
	}

	var input CreatePaymentInput

	// 🧪 Validate JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🧠 1. Get rent cycle
	cycles, err := database.Get("rent_cycles", "?id=eq."+input.RentCycleID)
	if err != nil || len(cycles) == 0 {
		c.JSON(404, gin.H{"error": "rent cycle not found"})
		return
	}

	cycle := cycles[0]
	tenantID := cycle["tenant_id"].(string)

	// 🧠 2. Get tenant
	tenants, err := database.Get("tenants", "?id=eq."+tenantID)
	if err != nil || len(tenants) == 0 {
		c.JSON(404, gin.H{"error": "tenant not found"})
		return
	}

	tenant := tenants[0]
	houseID := tenant["house_id"].(string)

	// 🧠 3. Ownership check (landlord)
	if role == "landlord" {

		users, _ := database.Get("users", "?email=eq."+email)
		userID := users[0]["id"].(string)

		landlords, _ := database.Get("landlords", "?user_id=eq."+userID)
		landlordID := landlords[0]["id"].(string)

		houses, _ := database.Get("houses", "?id=eq."+houseID+"&landlord_id=eq."+landlordID)

		if len(houses) == 0 {
			c.JSON(403, gin.H{"error": "unauthorized"})
			return
		}
	}

	// 🧠 (Optional) caretaker validation can be added here

	// 🧠 4. Insert payment
	data := map[string]interface{}{
		"rent_cycle_id": input.RentCycleID,
		"amount":        input.Amount,
		"method":        input.Method,
		"payment_date":  time.Now(),
	}

	err = database.Insert("payments", data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "payment recorded",
	})
}

func GetTenantBalance(c *gin.Context) {

	tenantID := c.Param("id")

	// 1. Get rent cycles
	cycles, err := database.Get("rent_cycles", "?tenant_id=eq."+tenantID)
	if err != nil || len(cycles) == 0 {
		c.JSON(404, gin.H{"error": "no rent cycles found"})
		return
	}

	var totalRent float64 = 0
	var totalPaid float64 = 0

	// 2. Loop cycles
	for _, cycle := range cycles {

		cycleID := cycle["id"].(string)
		rent := cycle["rent_amount"].(float64)

		totalRent += rent

		// 3. Get payments for this cycle
		payments, _ := database.Get("payments", "?rent_cycle_id=eq."+cycleID)

		for _, p := range payments {
			totalPaid += p["amount"].(float64)
		}
	}

	balance := totalRent - totalPaid

	c.JSON(200, gin.H{
		"total_rent": totalRent,
		"paid":       totalPaid,
		"balance":    balance,
		"in_arrears": balance > 0,
	})
}