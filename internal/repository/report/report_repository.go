package repository

import (
	"app/main/internal/repository"
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
)

type reportRepository struct {
	grpc proto.ReportUsersServiceClient
}

const reportRepositoryKey = "REPORT_GRPC_HOST"

func New() repository.Interface {
	return &reportRepository{}
}

func (r *reportRepository) Init() error {

	if r.grpc == nil {
		host := os.Getenv(reportRepositoryKey)
		if len(host) == 0 {
			return fmt.Errorf("report repository environment not found")
		}

		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("did not connect: %v", err)
		}

		r.grpc = proto.NewReportUsersServiceClient(conn)
	}
	return nil
}

func (r *reportRepository) Add(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.AddReportRequst); ok {
		return r.grpc.AddReport(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (r *reportRepository) Get(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetReportRequest); ok {
		return r.grpc.GetReport(ctx, val)
	}
	return nil, fmt.Errorf("invalid input parameter")
}

func (r *reportRepository) IsExist(req interface{}) (bool, error) {
	return true, nil
}

func (r *reportRepository) Verify(req interface{}) (interface{}, error) {
	return 1, nil
}
