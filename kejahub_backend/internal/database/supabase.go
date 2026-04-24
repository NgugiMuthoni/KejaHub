package database

import (
	"log"
	"os"
)

// Connect just validates environment variables
func Connect() {
	log.Println("🚀 Checking Supabase configuration...")

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	if url == "" {
		log.Fatal("❌ SUPABASE_URL is missing in .env")
	}

	if key == "" {
		log.Fatal("❌ SUPABASE_KEY is missing in .env")
	}

	log.Println("✅ Supabase configuration loaded successfully")
}
