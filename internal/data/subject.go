package data

import (
	"context"
	"edustate/internal/biz"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type subjectRepo struct {
	data *Data
	log  *log.Helper
}

func NewSubjectRepo(data *Data, logger log.Logger) biz.SubjectRepo {
	return &subjectRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *subjectRepo) GetByID(ctx context.Context, id int64) (*biz.Subject, error) {
	var subject biz.Subject
	err := r.data.db.WithContext(ctx).First(&subject, id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Errorf("GetByID err: %+v", err)
	}
	return &subject, nil
}
