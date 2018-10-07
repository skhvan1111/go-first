package redis

import (
	"github.com/go-redis/redis"
)

var client *redis.Client

func Init(addr string, password string) {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

func IncreaseCounter(keyName string) (int64, error) {
	counter, err := client.Incr(keyName).Result()
	return counter, err
}
