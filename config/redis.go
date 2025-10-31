package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis initializes the Redis connection
func InitRedis(redisURL string) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Failed to parse Redis URL, using default connection: %v", err)
		// Fallback to default connection
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	} else {
		RedisClient = redis.NewClient(opt)
	}

	// Test connection
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("WARNING: Failed to connect to Redis: %v", err)
		log.Println("Session management will fall back to in-memory storage")
		RedisClient = nil
		return
	}

	log.Println("Connected to Redis successfully")
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

// SetSession stores a session in Redis with expiration
func SetSession(token string, userID int, role string, expiration time.Duration) error {
	if RedisClient == nil {
		return nil // Silent fail for in-memory fallback
	}

	sessionData := map[string]interface{}{
		"user_id": userID,
		"role":    role,
	}

	return RedisClient.HSet(ctx, "session:"+token, sessionData).Err()
}

// GetSession retrieves a session from Redis
func GetSession(token string) (userID int, role string, exists bool, err error) {
	if RedisClient == nil {
		return 0, "", false, nil // Fall back to in-memory
	}

	result, err := RedisClient.HGetAll(ctx, "session:"+token).Result()
	if err != nil {
		return 0, "", false, err
	}

	if len(result) == 0 {
		return 0, "", false, nil
	}

	// Parse user_id
	var id int
	if userIDStr, ok := result["user_id"]; ok {
		_, err = fmt.Sscanf(userIDStr, "%d", &id)
		if err != nil {
			return 0, "", false, err
		}
	}

	role = result["role"]
	return id, role, true, nil
}

// DeleteSession removes a session from Redis
func DeleteSession(token string) error {
	if RedisClient == nil {
		return nil
	}

	return RedisClient.Del(ctx, "session:"+token).Err()
}
