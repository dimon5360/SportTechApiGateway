package repository

import (
	"app/main/internal/dto"
	"app/main/internal/repository"
	"context"
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

func (r *authRepository) Login(req *dto.RestLoginRequest) (*dto.RestLoginResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := r.grpc.LoginUser(ctx, dto.ConvertRest2GrpcLoginRequest(req))
	if err != nil {
		return nil, err
	}

	return dto.ConvertGrpc2RestLoginResponse(response), nil
}

func (r *authRepository) Register(req *dto.RestRegisterRequest) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcResponse, err := r.grpc.RegisterUser(ctx, dto.ConvertRest2GrpcRegisterRequest(req))
	if err != nil {
		return err
	}

	return dto.ConvertGrpc2RestRegisterResponse(grpcResponse)
}

func (r *authRepository) Refresh(req *dto.RestRefreshTokenRequest) (*dto.RestRefreshTokenResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := r.grpc.RefreshToken(ctx, dto.ConvertRest2GrpcRefreshRequest(req))
	if err != nil {
		return nil, err
	}
	return dto.ConvertGrpc2RestRefreshnResponse(response), nil
}
