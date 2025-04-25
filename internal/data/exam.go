package data

import (
	"context"
	"edustate/internal/biz"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type examRepo struct {
	data *Data
	log  *log.Helper
}

func NewExamRepo(data *Data, logger log.Logger) biz.ExamRepo {
	return &examRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (r *examRepo) GetByID(ctx context.Context, id int64) (*biz.Exam, error) {
	var exam biz.Exam
	err := r.data.db.WithContext(ctx).First(&exam, id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Errorf("examRepo.GetByID err: %+v", err)
	}
	return &exam, err
}
