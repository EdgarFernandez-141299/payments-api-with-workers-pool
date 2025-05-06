package db

import (
	"gitlab.com/clubhub.ai1/gommon/sql/common"
	"gitlab.com/clubhub.ai1/gommon/sql/gorm"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	gOrm "gorm.io/gorm"
)

const (
	maxIdle  = 1
	logQuery = true
	singular = true
)

func NewPostgresConnection() *gOrm.DB {
	postgresCfg := config.Config().Db

	return gorm.CreateDataBaseConnection(common.Connection{
		User:        postgresCfg.User,
		Password:    postgresCfg.Password,
		Host:        postgresCfg.Host,
		Port:        postgresCfg.Port,
		DBName:      postgresCfg.DB,
		DBProvider:  common.Postgresql,
		PoolSize:    postgresCfg.PoolSize,
		MaxIdleTime: maxIdle,
	}, common.AdditionalConfig{
		LogQuery:      logQuery,
		SingularTable: singular,
	})
}
