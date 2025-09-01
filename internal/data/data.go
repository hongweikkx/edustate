package data

import (
	"edustate/internal/conf"
	gormsql "edustate/pkg/gorm"
	redis_pkg "edustate/pkg/redis"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewExamRepo,
	NewScoreRepo,
	NewScoreItemRepo,
	NewStudentRepo,
	NewSubjectRepo,
)

// Data 是对所有数据库资源的统一封装
type Data struct {
	db    *gorm.DB
	redis *redis_pkg.Client
}

// NewData 创建 Data 并注入所有 repo 所需依赖
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, closeSqlF, err := gormsql.Init(logger, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	redis, closeRedisF, err := redis_pkg.Init(c)
	if err != nil {
		closeSqlF()
		return nil, nil, err
	}
	d := &Data{db: db, redis: redis}
	closeF := func() {
		closeSqlF()
		closeRedisF()
	}
	return d, closeF, nil
}
