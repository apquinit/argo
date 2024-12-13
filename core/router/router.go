package router

import (
	"github.com/gin-gonic/gin"
)

var routeRegistrations []func(*gin.Engine)

// InitRouter initializes the Gin router
func InitRouter() *gin.Engine {
	router := gin.Default()

	// Register all routes
	for _, register := range routeRegistrations {
		register(router)
	}

	return router
}

// RegisterRoutes adds a new route registration function
func RegisterRoutes(register func(*gin.Engine)) {
	routeRegistrations = append(routeRegistrations, register)
}
