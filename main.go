package main

import (
	"snowplow/schema"

	"github.com/gin-gonic/gin"
)

type client struct {
	Router *gin.Engine
}

var c client

func main() {
	c = client{
		Router: schema.InitRoutes(),
	}
	c.Router.Run("0.0.0.0:8080")
}
