package bloom

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v7"
)

func initMockRedis() *redis.Client {
	s, _ := miniredis.Run()
	return redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})
}
