package services

import (
	"kejahub-backend/internal/database"
	"sort"
)

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func ApplyPayment(cycle map[string]interface{}, paymentAmount float64) error {

	cycleID := cycle["id"].(string)
	remaining := paymentAmount

	// GET CHARGES (utilities only)
	charges, _ := database.GetAll(
		"charges",
		"?rent_cycle_id=eq."+cycleID,
	)

	// PRIORITY FIRST
	sort.Slice(charges, func(i, j int) bool {
		return charges[i]["priority"].(bool)
	})

	// APPLY TO UTILITIES
	for _, c := range charges {

		if remaining <= 0 {
			break
		}

		total := c["amount"].(float64)
		paid := c["paid_amount"].(float64)

		balance := total - paid
		if balance <= 0 {
			continue
		}

		apply := min(balance, remaining)

		database.Update("charges", c["id"].(string), map[string]interface{}{
			"paid_amount": paid + apply,
		})

		remaining -= apply
	}

	// APPLY TO RENT
	if remaining > 0 {

		rent := cycle["rent_amount"].(float64)
		paid := cycle["paid_amount"].(float64)

		balance := rent - paid

		if balance > 0 {

			apply := min(balance, remaining)

			database.Update("rent_cycles", cycleID, map[string]interface{}{
				"paid_amount": paid + apply,
			})
		}
	}

	return nil
}
