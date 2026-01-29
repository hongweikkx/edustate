package biz

import (
	"context"
	"edustate/internal/conf"
	"edustate/pkg/eino"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type AnalysisUseCase struct {
	scoreRepo ScoreRepo
	c         *conf.LLM
}

func NewAnalysisUseCase(scoreRepo ScoreRepo, c *conf.LLM) *AnalysisUseCase {
	return &AnalysisUseCase{
		scoreRepo: scoreRepo,
		c:         c,
	}
}

// Analyze 是核心业务逻辑：基于 studentID 查询成绩，并返回总结 + 建议
func (uc *AnalysisUseCase) Analyze(ctx context.Context, nlInputStr string) (string, []string, error) {
	studentID, err := eino.LLMClient.NLToArgs(ctx, nlInputStr)
	if err != nil {
		return "", nil, err
	}
	scores, err := uc.scoreRepo.GetByStudentID(ctx, cast.ToInt64(studentID))
	if err != nil {
		return "", nil, err
	}
	for _, score := range scores {
		if score.TotalScore < 60 {
			log.Context(ctx).Infof("学生 %d 在考试 %d 中的成绩为 %f，低于及格线", score.StudentID, score.ExamID, score.TotalScore)
			return "成绩分析完成", []string{}, nil
		}
	}
	return "成绩分析完成", []string{}, nil
}
