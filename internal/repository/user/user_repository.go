package repository

import (
	"app/main/internal/repository"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "proto/go"
)

type userRepository struct {
	grpc proto.UserServiceClient
}

const userRepositoryKey = "USER_SERVICE_HOST"

func New() repository.Interface {

	return &userRepository{}
}

func (r *userRepository) Init() error {

	if r.grpc == nil {
		host := os.Getenv(userRepositoryKey)
		if len(host) == 0 {
			log.Fatal("user repository environment not found")
		}

		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		r.grpc = proto.NewUserServiceClient(conn)
	}
	return nil
}

func (r *userRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.CreateUserRequest); ok {
		return r.grpc.CreateUser(ctx, val)
	}
	return nil, fmt.Errorf(repository.InvalidInputParameter)
}

func (r *userRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetUserRequest); ok {
		return r.grpc.GetUser(ctx, val)
	}
	return nil, fmt.Errorf(repository.InvalidInputParameter)
}

func (r *userRepository) Update(req interface{}) (interface{}, error) {
	return true, nil
}

func (r *userRepository) Delete(req interface{}) error {
	return nil
}
