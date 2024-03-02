package repository

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

type profileRepository struct {
	grpc proto.ProfileUsersServiceClient
}

func NewProfileRepository() *profileRepository {
	return &profileRepository{}
}

func (s *profileRepository) Init() error {

	conn, err := grpc.Dial(utils.Env().Value("PROFILE_GRPC_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.grpc = proto.NewProfileUsersServiceClient(conn)

	return nil
}

func (s *profileRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.CreateProfileRequst); ok {
		return s.grpc.CreateProfile(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (s *profileRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetProfileRequest); ok {
		return s.grpc.GetProfile(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}
