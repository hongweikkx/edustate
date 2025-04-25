package data

import (
	"context"
	"time"

	"edustate/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
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
	db *gorm.DB
}

// NewData 创建 Data 并注入所有 repo 所需依赖
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	gormLogger := NewGormLogger(logger, gormlog.Info)
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, nil, err
	}
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	d := &Data{db: db}
	cleanup := func() {
		_ = sqlDB.Close()
	}
	return d, cleanup, nil
}

type GormLogger struct {
	logger *log.Helper
	level  gormlog.LogLevel
}

func NewGormLogger(logger log.Logger, level gormlog.LogLevel) gormlog.Interface {
	return &GormLogger{
		logger: log.NewHelper(logger),
		level:  level,
	}
}

func (l *GormLogger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	return &GormLogger{
		logger: l.logger,
		level:  level,
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlog.Info {
		l.logger.WithContext(ctx).Infow(append([]interface{}{"msg", msg}, data...)...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlog.Warn {
		l.logger.WithContext(ctx).Warnw(append([]interface{}{"msg", msg}, data...)...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlog.Error {
		l.logger.WithContext(ctx).Errorw(append([]interface{}{"msg", msg}, data...)...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []interface{}{
		"duration", elapsed,
		"rows", rows,
		"sql", sql,
	}

	switch {
	case err != nil && l.level >= gormlog.Error:
		l.logger.WithContext(ctx).Errorw(append([]interface{}{"msg", "GORM Trace"}, append(fields, "error", err)...)...)
	case elapsed > 200*time.Millisecond && l.level >= gormlog.Warn:
		l.logger.WithContext(ctx).Warnw(append([]interface{}{"msg", "GORM Slow SQL"}, fields...)...)
	case l.level >= gormlog.Info:
		l.logger.WithContext(ctx).Infow(append([]interface{}{"msg", "GORM Info"}, fields...)...)
	}
}
