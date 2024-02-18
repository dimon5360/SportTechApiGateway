package grpc_service

import (
	"context"
	"log"
	"time"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserInfo struct {
	UserId      uint64
	AccessToken string
}

type AuthService struct {
	grpc proto.AuthUsersServiceClient

	Users map[uuid.UUID]UserInfo
}

func NewAuthService(host string) *AuthService {

	var s AuthService
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.grpc = proto.NewAuthUsersServiceClient(conn)

	s.Users = make(map[uuid.UUID]UserInfo)
	return &s
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
