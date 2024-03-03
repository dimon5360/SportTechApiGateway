package repository

import (
	"app/main/internal/repository"
	model "app/main/internal/repository/model"
	"app/main/pkg/utils"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

import (
	"context"
)

const (
	invalidRedisReq  = "invalid redis request"
	redisHostKey     = "REDIS_DB_HOST"
	redisPasswordKey = "REDIS_ADMIN_PASSWORD"
)

type redisRepository struct {
	client *redis.Client
}

var repo *redisRepository

func NewTokenRepository() repository.Interface {
	if repo == nil {
		repo = &redisRepository{
			client: nil,
		}
	}
	return repo
}

func (r *redisRepository) Init() error {

	if r.client == nil {
		env := utils.Env()
		ip, err := env.Value(redisHostKey)
		if err != nil {
			log.Fatal(err)
		}
		pass, err := env.Value(redisPasswordKey)
		if err != nil {
			log.Fatal(err)
		}

		opt := redis.Options{
			Addr:     ip,
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
