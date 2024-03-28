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

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
)

type userRepository struct {
	grpc proto.AuthUsersServiceClient
}

const userRepositoryKey = "USER_GRPC_HOST"

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

		r.grpc = proto.NewAuthUsersServiceClient(conn)
	}
	return nil
}

func (r *userRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetUserRequest); ok {
		return r.grpc.GetUser(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (r *userRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.CreateUserRequst); ok {
		return r.grpc.CreateUser(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (r *userRepository) IsExist(req interface{}) (bool, error) {
	return true, nil
}

func (r *userRepository) Verify(req interface{}) (interface{}, error) {
	return &proto.UserInfoResponse{
		Id: 1,
	}, nil
}

// func (s *userRepository) Auth(req *proto.AuthUserRequest) (*proto.UserInfoResponse, error) {

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	return s.grpc.AuthUser(ctx, req)
// }
