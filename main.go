package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type client struct {
	Router *gin.Engine
}

var c client
var redisClient *redis.Client

func init() {
	c = client{
		Router: initRoutes(),
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	InfoLogger.Printf(`successfully initialized :)`)
}

func main() {
	c.Router.Run("localhost:8080")
	InfoLogger.Println(`started router`)
}
