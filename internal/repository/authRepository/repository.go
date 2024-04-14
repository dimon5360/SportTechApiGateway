package authRepository

import (
	"app/main/internal/dto/models"
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "proto/go"
)

type Interface interface {
	Init() error
	Register(*models.RestRegisterRequest) error
	Login(*models.RestLoginRequest) (*models.RestLoginResponse, error)
	RefreshToken(*models.RestRefreshTokenRequest) (*models.RestRefreshTokenResponse, error)
}

type authRepository struct {
	grpc proto.AuthServiceClient
}

const (
	authRepositoryKey = "AUTH_SERVICE_HOST"
)

func New() Interface {
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

func (r *authRepository) Register(req *models.RestRegisterRequest) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	grpcResponse, err := r.grpc.RegisterUser(ctx, models.ConvertRest2GrpcRegisterRequest(req))
	if err != nil {
		return err
	}

	return models.ConvertGrpc2RestRegisterResponse(grpcResponse)
}

func (r *authRepository) Login(req *models.RestLoginRequest) (*models.RestLoginResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := r.grpc.LoginUser(ctx, models.ConvertRest2GrpcLoginRequest(req))
	if err != nil {
		return nil, err
	}

	return models.ConvertGrpc2RestLoginResponse(response), nil
}

func (r *authRepository) RefreshToken(req *models.RestRefreshTokenRequest) (*models.RestRefreshTokenResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := r.grpc.RefreshToken(ctx, models.ConvertRest2GrpcRefreshRequest(req))
	if err != nil {
		return nil, err
	}
	return models.ConvertGrpc2RestRefreshnResponse(response), nil
}

func (r *authRepository) Delete(req interface{}) error {
	return nil
}
