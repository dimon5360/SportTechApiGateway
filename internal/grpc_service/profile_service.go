package grpc_service

import (
	"context"
	"log"
	"time"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProfileService struct {
	grpc proto.ProfileUsersServiceClient
}

func NewProfileService(host string) *ProfileService {

	var s ProfileService
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.grpc = proto.NewProfileUsersServiceClient(conn)
	return &s
}

func (s *ProfileService) CreateProfile(req *proto.CreateProfileRequst) (*proto.UserProfileResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.CreateProfile(ctx, req)
}

func (s *ProfileService) GetProfile(req *proto.GetProfileRequest) (*proto.UserProfileResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.GetProfile(ctx, req)
}
