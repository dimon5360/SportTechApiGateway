package repository

import (
	"app/main/internal/repository"
	model "app/main/internal/repository/model"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

import (
	"context"
)

const (
	invalidRedisReq  = "invalid redis request"
	redisHostKey     = "REDIS_HOST"
	redisPasswordKey = "REDIS_ADMIN_PASSWORD"
)

type redisRepository struct {
	client *redis.Client
}

var repo *redisRepository

func New() repository.Interface {
	if repo == nil {
		repo = &redisRepository{
			client: nil,
		}
	}
	return repo
}

func (r *redisRepository) Init() error {

	if r.client == nil {
		host := os.Getenv(redisHostKey)
		if len(host) == 0 {
			log.Fatal("redis host not found")
		}
		pass := os.Getenv(redisPasswordKey)
		if len(pass) == 0 {
			log.Fatal("redis password not found")
		}

		opt := redis.Options{
			Addr:     host,
			Password: pass, // no password set
			DB:       0,    // use default DB
		}
		r.client = redis.NewClient(&opt)
	}
	return nil
}

func (r *redisRepository) Add(req interface{}) (interface{}, error) {
	ctx := context.Background()

	data, ok := req.(*model.RedisRequestModel)
	if !ok {
		return nil, fmt.Errorf(invalidRedisReq)
	}

	err := r.client.Set(ctx, data.Key, data.Value, data.Expire).Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nil, nil
}

func (r *redisRepository) Get(req interface{}) (interface{}, error) {
	ctx := context.Background()

	data, ok := req.(*model.RedisRequestModel)
	if !ok {
		return nil, fmt.Errorf(invalidRedisReq)
	}

	val, err := r.client.Get(ctx, data.Key).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return []byte(val), nil
}

func (r *redisRepository) IsExist(req interface{}) (bool, error) {
	return true, nil
}

func (r *redisRepository) Verify(req interface{}) (interface{}, error) {
	return 1, nil
}
