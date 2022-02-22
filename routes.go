package main

import (
	"github.com/gin-gonic/gin"
)

// initRoutes initializes the API routes
func initRoutes() *gin.Engine {
	router := gin.Default()
	// /schema route group
	routeGroup := router.Group("/schema")
	{
		// create schema
		routeGroup.POST(schemaEndpoint, addSchema)
		// get schema
		routeGroup.GET(schemaEndpoint, getSchema)

		// nested /validate route group
		validateGroup := routeGroup.Group("/validate")
		{
			// validate JSON using schema
			validateGroup.POST(schemaEndpoint, validateJSON)
		}
	}

	return router
}
