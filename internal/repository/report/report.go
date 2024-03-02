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

type reportRepository struct {
	grpc proto.ReportUsersServiceClient
}

func NewReportRepository() *reportRepository {
	return &reportRepository{}
}

func (s *reportRepository) Init() error {

	conn, err := grpc.Dial(utils.Env().Value("REPORT_GRPC_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.grpc = proto.NewReportUsersServiceClient(conn)

	return nil
}

func (s *reportRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.AddReportRequst); ok {
		return s.grpc.AddReport(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (s *reportRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetReportRequest); ok {
		return s.grpc.GetReport(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}
