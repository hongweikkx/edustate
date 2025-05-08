package service

import (
	"context"
	"edustate/internal/biz"
	"github.com/go-kratos/kratos/v2/log"

	pb "edustate/api/edustate/v1"
)

type AnalysisService struct {
	pb.UnimplementedAnalysisServer
	uc  *biz.AnalysisUsecase
	log *log.Helper
}

func NewAnalysisService(uc *biz.AnalysisUsecase, logger log.Logger) *AnalysisService {
	return &AnalysisService{uc: uc, log: log.NewHelper(logger)}
}

func (s *AnalysisService) Analyze(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeReply, error) {
	summary, suggestions, err := s.uc.Analyze(ctx, req.GetStudentNlInput())
	if err != nil {
		return nil, err
	}
	return &pb.AnalyzeReply{
		ResultSummary: summary,
		Suggestions:   suggestions,
	}, nil
}
