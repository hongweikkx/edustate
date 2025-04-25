package data

import (
	"context"
	"edustate/internal/biz"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

func NewStudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *studentRepo) GetByID(ctx context.Context, id int64) (*biz.Student, error) {
	var student biz.Student
	err := r.data.db.WithContext(ctx).First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepo) GetByStudentNumber(ctx context.Context, sn string) (*biz.Student, error) {
	var student biz.Student
	err := r.data.db.WithContext(ctx).Where("student_number = ?", sn).First(&student).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("studentRepo.GetByStudentNumber err: %+v", err)
	}
	return &student, err
}
