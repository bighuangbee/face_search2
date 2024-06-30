package storage

import (
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type MysqlStorage struct {
	db *gorm.DB
}

func NewMysqlStorage(config *conf.Data, logger log.Logger) (*MysqlStorage, error) {
	db, err := NewDB(config, logger)
	db.AutoMigrate(&RegisteInfo{})
	return &MysqlStorage{db: db}, err
}

func (s *MysqlStorage) Update(key string, value *RegisteInfo) error {
	return s.db.Where("filename = ?", key).Save(value).Error
}

func (s *MysqlStorage) Read(key string) (value *RegisteInfo, ok bool) {
	value = &RegisteInfo{}
	if err := s.db.Where("filename = ?", key).First(value).Error; err != nil {
		return nil, false
	}
	return value, true
}

func (s *MysqlStorage) ReadBatch() (values []*RegisteInfo, err error) {
	if err := s.db.Find(&values).Error; err != nil {
		return nil, err
	}
	return
}

func (s *MysqlStorage) Delete(key string) error {
	return s.db.Where("filename = ?", key).Delete(&RegisteInfo{}).Error
}

func (s *MysqlStorage) Count() (count int64) {
	s.db.Count(&count)
	return
}

func (s *MysqlStorage) DeleteExpired(effectiveTime time.Duration) (values []*RegisteInfo, err error) {
	cutoffTime := time.Now().Add(-effectiveTime)
	if err := s.db.Where("shoot_time < ?", cutoffTime).Find(&values).Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("shoot_time < ?", cutoffTime).Delete(&RegisteInfo{}).Error; err != nil {
		return nil, err
	}

	return values, nil
}

func NewDB(config *conf.Data, logger log.Logger) (*gorm.DB, error) {
	logs := log.NewHelper(log.With(logger, "module", "receive-service/data/gorm"))

	db, err := gorm.Open(mysql.Open(config.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		logs.Fatalf("failed opening connection to mysql: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(config.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(config.Database.MaxOpenConns))

	return db, nil
}
