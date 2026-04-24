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

func MpesaCallback(c *gin.Context) {

	var payload map[string]interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔍 DEBUG: print full callback
	// (you can remove later)
	println("MPESA CALLBACK RECEIVED")

	// 🧠 Extract data safely
	body, ok := payload["Body"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "invalid format"})
		return
	}

	stkCallback, ok := body["stkCallback"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "invalid callback"})
		return
	}

	resultCode := int(stkCallback["ResultCode"].(float64))

	// ❌ failed payment
	if resultCode != 0 {
		c.JSON(200, gin.H{"message": "payment failed"})
		return
	}

	callbackMetadata := stkCallback["CallbackMetadata"].(map[string]interface{})
	items := callbackMetadata["Item"].([]interface{})

	var amount float64
	var phone string

	for _, item := range items {
		entry := item.(map[string]interface{})

		if entry["Name"] == "Amount" {
			amount = entry["Value"].(float64)
		}

		if entry["Name"] == "PhoneNumber" {
			phone = entry["Value"].(string)
		}
	}

	// ⚠️ For now: manual mapping (we’ll improve later)
	// You should map phone → tenant → rent_cycle

	data := map[string]interface{}{
		"amount": amount,
		"method": "mpesa",
	}

	database.Insert("payments", data)

	c.JSON(200, gin.H{"message": "callback processed", "phone": phone})
}

func STKPushRequest(c *gin.Context) {

	var input STKInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := services.STKPush(input.Phone, input.Amount)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 🔥 STORE CHECKOUT REQUEST ID (IMPORTANT FIX)
	database.Insert("payments", map[string]interface{}{
		"phone":               input.Phone,
		"amount":              input.Amount,
		"checkout_request_id": result.CheckoutRequestID,
		"status":              "pending",
		"method":              "mpesa",
	})

	c.JSON(200, gin.H{
		"checkout_request_id": result.CheckoutRequestID,
		"message":             result.CustomerMessage,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
