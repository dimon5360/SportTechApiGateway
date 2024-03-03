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

const invalidRedisReq = "invalid redis request"

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository() repository.Interface {
	return &redisRepository{
		client: nil,
	}
}

func (r *redisRepository) Init() error {

	ip, err := utils.Env().Value("REDIS_DB_HOST")
	if err != nil {
		log.Fatal(err)
	}
	pass, err := utils.Env().Value("REDIS_HOST_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}

	opt := redis.Options{
		Addr:     ip,
		Password: pass, // no password set
		DB:       0,    // use default DB
	}

	r.client = redis.NewClient(&opt)
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
