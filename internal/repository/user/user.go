package grpc_service

import (
	"app/main/pkg/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userRepository struct {
	grpc proto.AuthUsersServiceClient
}

func NewUserRepository() *userRepository {

	return &userRepository{}
}

func (s *userRepository) Init() error {

	conn, err := grpc.Dial(utils.Env().Value("USER_GRPC_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.grpc = proto.NewAuthUsersServiceClient(conn)

	return nil
}

func (s *userRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetUserRequest); ok {
		return s.grpc.GetUser(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (s *userRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.CreateUserRequst); ok {
		return s.grpc.CreateUser(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

// func (s *userRepository) Auth(req *proto.AuthUserRequest) (*proto.UserInfoResponse, error) {

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	return s.grpc.AuthUser(ctx, req)
// }