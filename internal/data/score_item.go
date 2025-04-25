package data

import (
	"context"
	"edustate/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type scoreItemRepo struct {
	data *Data
	log  *log.Helper
}

func NewScoreItemRepo(data *Data, logger log.Logger) biz.ScoreItemRepo {
	return &scoreItemRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *scoreItemRepo) ListByScoreID(ctx context.Context, scoreID int64) ([]*biz.ScoreItem, error) {
	var items []*biz.ScoreItem
	err := r.data.db.WithContext(ctx).Where("score_id = ?", scoreID).Find(&items).Error
	if err != nil {
		r.log.Errorf("query scoreItemRepo.ListByScoreID err: %+v", err)
		return nil, err
	}
	return items, nil
}
