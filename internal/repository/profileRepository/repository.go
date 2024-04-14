package profileRepository

import (
	"app/main/internal/dto/constants"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "proto/go"
)

type Interface interface {
	Init() error
	Create(interface{}) (interface{}, error)
	Read(interface{}) (interface{}, error)
	Update(interface{}) (interface{}, error)
	Delete(interface{}) error
}

type profileRepository struct {
	grpc proto.UserServiceClient
}

const profileRepositoryKey = "USER_SERVICE_HOST"

func New() Interface {
	return &profileRepository{}
}

func (r *profileRepository) Init() error {

	if r.grpc == nil {
		host := os.Getenv(profileRepositoryKey)
		if len(host) == 0 {
			log.Fatal("profile repository environment not found")
		}

		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		r.grpc = proto.NewUserServiceClient(conn)
	}
	return nil
}

func (r *profileRepository) Create(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.CreateProfileRequest); ok {
		return r.grpc.CreateProfile(ctx, val)
	}
	return nil, fmt.Errorf(constants.InvalidInputParameter)
}

func (r *profileRepository) Read(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetProfileRequest); ok {
		return r.grpc.GetProfile(ctx, val)
	}
	return nil, fmt.Errorf(constants.InvalidInputParameter)
}

func (r *profileRepository) Update(req interface{}) (interface{}, error) {
	return true, nil
}

func (r *profileRepository) Delete(req interface{}) error {
	return nil
}
