package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ComputedDetail(c *gin.Context) {
	id := c.Query("id")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":                    id,
		"replace":               "replace",
		"replace_if_configured": "replace_if_configured",
		"use_state_for_unknown": "use_state",
		"list_optional":         []string{"list_optional"},
	})
}
