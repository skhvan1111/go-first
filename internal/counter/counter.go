package counter

import (
	"github.com/skhvan1111/go-first/internal/redis"
)

func Increase() (int64, error) {
	return redis.IncreaseCounter("NumberOfVisits")
}
