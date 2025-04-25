package data

import (
	"context"
	"edustate/internal/biz"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type scoreRepo struct {
	data *Data
	log  *log.Helper
}

func NewScoreRepo(data *Data, logger log.Logger) biz.ScoreRepo {
	return &scoreRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *scoreRepo) GetByExamSubjectStudent(ctx context.Context, examID, subjectID, studentID int64) (*biz.Score, error) {
	var s biz.Score
	err := r.data.db.WithContext(ctx).Where("exam_id = ? AND subject_id = ? AND student_id = ?", examID, subjectID, studentID).First(&s).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Errorf("scoreRepo.GetByExamSubjectStudent err: %+v", err)
	}
	return &s, err
}

func (r *scoreRepo) GetByStudentID(ctx context.Context, studentID int64) ([]*biz.Score, error) {
	var scores []*biz.Score
	err := r.data.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&scores).Error
	if err != nil {
		r.log.Errorf("scoreRepo.GetByStudentID err: %+v", err)
		return nil, err
	}
	return scores, nil
}
