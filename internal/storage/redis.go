package storage

import (
	"app/main/utils"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	client *redis.Client
}

func InitRedis() *RedisConnection {

	var conn RedisConnection

	ip := utils.Env().Value("REDIS_HOST")
	pass := utils.Env().Value("REDIS_HOST_PASSWORD")

	opt := redis.Options{
		Addr:     ip,
		Password: pass, // no password set
		DB:       0,    // use default DB
	}

	conn.client = redis.NewClient(&opt)

	return &conn
}

func (c *RedisConnection) TestConnect() {
	ctx := context.Background()

	err := c.client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := c.client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
}
