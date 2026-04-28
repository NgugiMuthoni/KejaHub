package routes

import (
	"fmt"
	"kejahub-backend/internal/services"
	"time"
	"kejahub-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func MpesaCallback(c *gin.Context) {

	var payload map[string]interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body := payload["Body"].(map[string]interface{})
	stk := body["stkCallback"].(map[string]interface{})

	resultCode := int(stk["ResultCode"].(float64))
	checkoutID := stk["CheckoutRequestID"].(string)

	// ❌ FAILED PAYMENT
	if resultCode != 0 {
		c.JSON(200, gin.H{"message": "failed"})
		return
	}

	// ✅ SUCCESS
	meta := stk["CallbackMetadata"].(map[string]interface{})
	items := meta["Item"].([]interface{})

	var amount float64
	var phone string

	for _, item := range items {
		entry := item.(map[string]interface{})

		if entry["Name"] == "Amount" {
			amount = entry["Value"].(float64)
		}

		if entry["Name"] == "PhoneNumber" {
			phone = fmt.Sprintf("%.0f", entry["Value"])
		}
	}

	// 🔍 FIND TENANT
	tenant, err := database.GetOne("tenants", "?phone=eq."+phone)
	if err != nil {
		c.JSON(200, gin.H{"message": "tenant not found"})
		return
	}

	tenantID := tenant["id"].(string)
	houseID := tenant["house_id"].(string)

	// 📅 ENSURE RENT CYCLE
	cycleID, _ := services.EnsureRentCycle(tenant)

	// 💾 STORE PAYMENT
	database.Insert("payments", map[string]interface{}{
		"tenant_id":     tenantID,
		"house_id":      houseID,
		"rent_cycle_id": cycleID,
		"amount":        amount,
		"method":        "mpesa",
		"payment_date":  time.Now(),
	})

	// 🔥 APPLY PAYMENT
	cycle, _ := database.GetOne("rent_cycles", "?id=eq."+cycleID)
	services.ApplyPayment(cycle, amount)

	c.JSON(200, gin.H{"message": "processed"})
}
