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

type authRepository struct {
	grpc proto.AuthServiceClient
}

const (
	authRepositoryKey = "AUTH_SERVICE_HOST"
)

func New() repository.AuthInterface {

	return &authRepository{}
}

func (r *authRepository) Init() error {

	if r.grpc == nil {
		host := os.Getenv(authRepositoryKey)
		if len(host) == 0 {
			log.Fatal("auth repository environment not found")
		}

		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		r.grpc = proto.NewAuthServiceClient(conn)
	}
	return nil
}

func (r *authRepository) Login(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return r.grpc.LoginUser(ctx, &proto.LoginUserRequest{
		Email:    "admin@test.com",
		Password: "admin1234",
	})

	// if val, ok := req.(*proto.LoginUserRequest); ok {
	// 	return r.grpc.LoginUser(ctx, val)
	// }
	// return nil, fmt.Errorf(repository.InvalidInputParameter)
}

func (r *authRepository) Register(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.RegisterUserRequest); ok {
		return r.grpc.RegisterUser(ctx, val)
	}
	return nil, fmt.Errorf(repository.InvalidInputParameter)
}

func (r *authRepository) Refresh(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.RefreshTokenRequest); ok {
		return r.grpc.RefreshToken(ctx, val)
	}
	return nil, fmt.Errorf(repository.InvalidInputParameter)
}
