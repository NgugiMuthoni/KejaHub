package routes

import (
	"kejahub-backend/internal/database"
	"kejahub-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type STKInput struct {
	Phone  string `json:"phone"`
	Amount int    `json:"amount"`
}

// 🚀 STK PUSH REQUEST
func STKPushRequest(c *gin.Context) {

	var input STKInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call STK service
	result, err := services.STKPush(input.Phone, input.Amount)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 💾 Store initial transaction (PENDING)
	database.Insert("payments", map[string]interface{}{
		"phone":               input.Phone,
		"amount":              input.Amount,
		"checkout_request_id": result.CheckoutRequestID,
		"status":              "pending",
		"method":              "mpesa",
	})

	// Return response immediately (IMPORTANT for STK)
	c.JSON(200, gin.H{
		"checkout_request_id": result.CheckoutRequestID,
		"message":             result.CustomerMessage,
	})
}

// 📩 MPESA CALLBACK
func MpesaCallback(c *gin.Context) {

	var payload map[string]interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body, ok := payload["Body"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "invalid callback body"})
		return
	}

	stkCallback, ok := body["stkCallback"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "invalid stk callback"})
		return
	}

	resultCode, ok := stkCallback["ResultCode"].(float64)
	if !ok {
		c.JSON(400, gin.H{"error": "invalid result code"})
		return
	}

	checkoutID, _ := stkCallback["CheckoutRequestID"].(string)

	// ❌ FAILED PAYMENT
	if int(resultCode) != 0 {
		database.Insert("payments", map[string]interface{}{
			"checkout_request_id": checkoutID,
			"status":              "failed",
			"method":              "mpesa",
		})

		c.JSON(200, gin.H{"message": "payment failed"})
		return
	}

	// ✅ SUCCESS CASE
	callbackMetadata, ok := stkCallback["CallbackMetadata"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "missing metadata"})
		return
	}

	items, ok := callbackMetadata["Item"].([]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "invalid items"})
		return
	}

	var amount float64
	var phone string

	for _, item := range items {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if entry["Name"] == "Amount" {
			amount, _ = entry["Value"].(float64)
		}

		if entry["Name"] == "PhoneNumber" {
			phone, _ = entry["Value"].(string)
		}
	}

	// 💾 Store successful payment
	database.Insert("payments", map[string]interface{}{
		"phone":               phone,
		"amount":              amount,
		"checkout_request_id": checkoutID,
		"status":              "paid",
		"method":              "mpesa",
	})

	c.JSON(200, gin.H{
		"message": "callback processed",
		"phone":   phone,
	})
}
