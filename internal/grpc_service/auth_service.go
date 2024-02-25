package grpc_service

import (
	"context"
	"log"
	"time"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	grpc proto.AuthUsersServiceClient

	isInitialized bool
}

var authService AuthService

func NewAuthService(host string) {

	if !authService.isInitialized {
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		authService.grpc = proto.NewAuthUsersServiceClient(conn)
		authService.isInitialized = true
	}
}

func AuthServiceInstance() *AuthService {

	if !authService.isInitialized {
		return nil
	}

	return &authService
}

func (s *AuthService) GetUser(req *proto.GetUserRequest) (*proto.UserInfoResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.GetUser(ctx, req)
}

func (s *AuthService) Auth(req *proto.AuthUserRequest) (*proto.UserInfoResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.AuthUser(ctx, req)
}

func (s *AuthService) Register(req *proto.CreateUserRequst) (*proto.UserInfoResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.CreateUser(ctx, req)
}
