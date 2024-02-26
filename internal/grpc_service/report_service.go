package grpc_service

import (
	"context"
	"log"
	"time"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReportService struct {
	grpc proto.ReportUsersServiceClient

	isInitialized bool
}

var reportService ReportService

func NewReportService(host string) {

	if !reportService.isInitialized {
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		reportService.grpc = proto.NewReportUsersServiceClient(conn)
		reportService.isInitialized = true
	}
}

func ReportServiceInstance() *ReportService {

	if !reportService.isInitialized {
		return nil
	}

	return &reportService
}

func (s *ReportService) CreateReport(req *proto.AddReportRequst) (*proto.ReportResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.AddReport(ctx, req)
}

func (s *ReportService) GetReport(req *proto.GetReportRequest) (*proto.ReportResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.grpc.GetReport(ctx, req)
}
