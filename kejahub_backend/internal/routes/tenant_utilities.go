package routes

import (
	"github.com/gin-gonic/gin"
	"kejahub-backend/internal/database"
)

func GetTenantUtilities(c *gin.Context) {

	tenantID := c.Param("id")

	cycles, _ := database.Get("rent_cycles", "?tenant_id=eq."+tenantID)

	var result []map[string]interface{}

	for _, cycle := range cycles {

		cycleID := cycle["id"].(string)

		utils, _ := database.Get("utilities", "?rent_cycle_id=eq."+cycleID)

		result = append(result, map[string]interface{}{
			"month":     cycle["month"],
			"utilities": utils,
		})
	}

	c.JSON(200, result)
}