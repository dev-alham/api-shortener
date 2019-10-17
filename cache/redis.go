package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"time"
)

var redisClient *redis.Client

func RedisInit() {
	redis_db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		PoolSize:   100,
		MaxRetries: 2,
		Password:   os.Getenv("REDIS_PASS"),
		DB:         redis_db,
	})

	ping, err := redisClient.Ping().Result()
	if err == nil && len(ping) > 0 {
		println("Connected to Redis")
	} else {
		println("Redis Connection Failed")
	}
}

func SetValueWithTTL(dir string, key string, value interface{}, ttl int) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	if dir != "" {
		key = fmt.Sprintf("%s:%s", dir, key)
	}
	err := redisClient.Set(key, string(serializedValue), time.Duration(ttl)*time.Minute).Err()
	return true, err
}

func GetValue(dir string, key string) (interface{}, error) {
	var deserializedValue interface{}
	if dir != "" {
		key = fmt.Sprintf("%s:%s", dir, key)
	}
	serializedValue, err := redisClient.Get(key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

func SetValue(dir string, key string, value interface{}) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	if dir != "" {
		key = fmt.Sprintf("%s:%s", dir, key)
	}
	err := redisClient.Set(key, string(serializedValue), 0).Err()
	return true, err
}

func RPush(key string, valueList []string) (bool, error) {
	err := redisClient.RPush(key, valueList).Err()
	return true, err
}

func RpushWithTTL(key string, valueList []string, ttl int) (bool, error) {
	err := redisClient.RPush(key, valueList, ttl).Err()
	return true, err
}
func LRange(key string) (bool, error) {
	err := redisClient.LRange(key, 0, -1).Err()
	return true, err
}
func ListLen(key string) int64 {
	return redisClient.LLen(key).Val()
}

func Publish(channel string, message string) {
	redisClient.Publish(channel, message)
}

func GetKeysByPattern(pattern string) []string {
	return redisClient.Keys(pattern).Val()
}

func Increment(key string) int64 {
	return redisClient.Incr(key).Val()
}

func DelKey(dir string, key string) error {
	if dir != "" {
		key = fmt.Sprintf("%s:%s", dir, key)
	}
	return redisClient.Del(key).Err()
}

func FlushRedis() error {
	return redisClient.FlushAll().Err()
}
