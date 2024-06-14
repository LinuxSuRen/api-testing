package writer_templates

import (
	"context"
	"log"
)

type ReportServer struct {
	UnimplementedReportWriterServer
}

func (s *ReportServer) SendReportResult(ctx context.Context, req *ReportResultRepeated) (*Empty, error) {
	// print received data
	for _, result := range req.Data {
		log.Printf("Received report: %+v", result)
	}
	return &Empty{}, nil
}
