package services

import (
	"fmt"
	"time"

	"kejahub-backend/internal/database"
)

// Ensure rent cycle exists
func EnsureRentCycle(tenant map[string]interface{}) (string, error) {

	tenantID := tenant["id"].(string)
	rent := tenant["rent_amount"].(float64)

	month := time.Now().Format("2006-01")

	// Check if exists
	cycle, err := database.GetOne(
		"rent_cycles",
		"?tenant_id=eq."+tenantID+"&month=eq."+month,
	)

	if err == nil {
		return cycle["id"].(string), nil
	}

	// Create new cycle
	err = database.Insert("rent_cycles", map[string]interface{}{
		"tenant_id":   tenantID,
		"month":       month,
		"rent_amount": rent,
		"due_date":    time.Now().AddDate(0, 0, 5),
	})

	if err != nil {
		return "", err
	}

	newCycle, err := database.GetOne(
		"rent_cycles",
		"?tenant_id=eq."+tenantID+"&month=eq."+month,
	)

	if err != nil {
		return "", fmt.Errorf("failed to fetch created cycle")
	}

	return newCycle["id"].(string), nil
}
