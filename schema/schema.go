package schema

import (
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	if ping := redisClient.Ping(); ping.Err() == nil {
		InfoLogger.Printf(`successfully connected to Redis`)
	}
}
