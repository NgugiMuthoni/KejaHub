package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"kejahub-backend/internal/database"
	"kejahub-backend/internal/middleware"
	"kejahub-backend/internal/routes"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found")
	}

	log.Println("MAIN STARTED")

	// Connect to Supabase
	database.Connect()
	database.Connect()

	// Create router
	router := gin.Default()

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", routes.Register)
		auth.POST("/login", routes.Login)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// House routes
	router.POST("/houses", middleware.AuthMiddleware(), routes.CreateHouse)
	// Caretaker routes
	router.POST("/caretakers", middleware.AuthMiddleware(), routes.CreateCaretaker)
	// User routes
	router.GET("/profile", middleware.AuthMiddleware(), routes.GetProfile)
	// Payment routes
	router.POST("/payments", middleware.AuthMiddleware(), routes.CreatePayment)
	router.GET("/tenants/:id/balance", middleware.AuthMiddleware(), routes.GetTenantBalance)
	// Dashboard routes
	router.GET("/dashboard", middleware.AuthMiddleware(), routes.LandlordDashboard)
	//tenant dashboard
	router.POST("/tenants", middleware.AuthMiddleware(), routes.CreateTenant)
	// Rent routes
	router.POST("/rent/generate", routes.GenerateMonthlyRent)
	//mpesa routes
	router.POST("/mpesa/callback", routes.MpesaCallback)
	// Utility routes
	router.POST("/utilities", middleware.AuthMiddleware(), routes.AddUtility)
	// Check priority utilities
	router.GET("/utilities/:id/check", middleware.AuthMiddleware(), routes.CheckPriorityUtilities)
	// Get tenant utilities
	router.GET("/tenants/:id/utilities", middleware.AuthMiddleware(), routes.GetTenantUtilities)
	router.POST("/mpesa/stk", middleware.AuthMiddleware(), routes.STKPushRequest)

	
	// Start server
	log.Println("Server running on :8080")
	router.Run(":8080")
}
