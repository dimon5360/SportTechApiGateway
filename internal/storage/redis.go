package storage

import (
	"app/main/utils"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	client *redis.Client
}

var conn RedisConnection

func InitRedis() {

	ip := utils.Env().Value("REDIS_HOST")
	pass := utils.Env().Value("REDIS_HOST_PASSWORD")

	opt := redis.Options{
		Addr:     ip,
		Password: pass, // no password set
		DB:       0,    // use default DB
	}

	conn.client = redis.NewClient(&opt)
}

func Redis() *RedisConnection {
	return &conn
}

func (c *RedisConnection) Store(key string, value []byte, expire time.Duration) {
	ctx := context.Background()

	err := c.client.Set(ctx, key, value, expire).Err()
	if err != nil {
		log.Println(err)
	}
}

func (c *RedisConnection) Get(key string) []byte {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		log.Println(err)
	}
	return []byte(val)
}
