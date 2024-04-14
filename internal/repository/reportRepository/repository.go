package reportRepository

import (
	"app/main/internal/dto/constants"
	"context"
	"fmt"
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

type reportRepository struct {
	grpc proto.ReportServiceClient
}

const reportRepositoryKey = "REPORT_SERVICE_HOST"

func New() Interface {
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

		r.grpc = proto.NewReportServiceClient(conn)
	}
	return nil
}

func (r *reportRepository) Create(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.AddReportRequst); ok {
		return r.grpc.AddReport(ctx, val)
	}
	return nil, fmt.Errorf(constants.InvalidInputParameter)
}

func (r *reportRepository) Read(req interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if val, ok := req.(*proto.GetReportRequest); ok {
		return r.grpc.GetReport(ctx, val)
	}
	return nil, fmt.Errorf(constants.InvalidInputParameter)
}

func (r *reportRepository) Update(req interface{}) (interface{}, error) {
	return true, nil
}

func (r *reportRepository) Delete(req interface{}) error {
	return nil
}
