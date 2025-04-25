package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AnalysisUsecase struct {
	scoreRepo ScoreRepo
	log       *log.Helper
}

func NewAnalysisUsecase(scoreRepo ScoreRepo, logger log.Logger) *AnalysisUsecase {
	return &AnalysisUsecase{
		scoreRepo: scoreRepo,
		log:       log.NewHelper(logger),
	}
}

// Analyze 是核心业务逻辑：基于 studentID 查询成绩，并返回总结 + 建议
func (uc *AnalysisUsecase) Analyze(ctx context.Context, studentID int64) (string, []string, error) {
	scores, err := uc.scoreRepo.GetByStudentID(ctx, studentID)
	if err != nil {
		return "", nil, err
	}
	for _, score := range scores {
		if score.TotalScore < 60 {
			uc.log.Infof("学生 %d 在考试 %d 中的成绩为 %f，低于及格线", score.StudentID, score.ExamID, score.TotalScore)
			return "成绩分析完成", []string{}, nil
		}
	}
	return "成绩分析完成", []string{}, nil
}
