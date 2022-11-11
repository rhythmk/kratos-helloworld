package data

import (
	"kratos-shop/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGormDB, NewGreeterRepo, NewStudentRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

func NewGormDB(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, cleanup, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, cleanup, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetConnMaxLifetime(time.Second * 25)
	return &Data{db: db}, cleanup, err
}

// // NewData .
// func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
// 	cleanup := func() {
// 		log.NewHelper(logger).Info("closing the data resources")
// 	}
// 	return &Data{}, cleanup, nil
// }
