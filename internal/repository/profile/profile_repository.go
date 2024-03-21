package repository

import (
	"app/main/internal/repository"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
)

type profileRepository struct {
	grpc proto.ProfileUsersServiceClient
}

const profileRepositoryKey = "PROFILE_GRPC_HOST"

func New() repository.Interface {
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
		r.grpc = proto.NewProfileUsersServiceClient(conn)
	}
	return nil
}

func (r *profileRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.CreateProfileRequst); ok {
		return r.grpc.CreateProfile(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (r *profileRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetProfileRequest); ok {
		return r.grpc.GetProfile(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (r *profileRepository) IsExist(req interface{}) (bool, error) {
	return true, nil
}

func (r *profileRepository) Verify(req interface{}) (interface{}, error) {
	return 1, nil
}
