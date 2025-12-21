package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-looklook/app/travel/model"
	"go-zero-looklook/app/travel/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis

	HomestayModel         model.HomestayModel
	HomestayBusinessModel model.HomestayBusinessModel
	HomestayActivityModel model.HomestayActivityModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                c,
		RedisClient:           redis.MustNewRedis(c.Redis.RedisConf),
		HomestayModel:         model.NewHomestayModel(sqlConn, c.Cache),
		HomestayActivityModel: model.NewHomestayActivityModel(sqlConn, c.Cache),
		HomestayBusinessModel: model.NewHomestayBusinessModel(sqlConn, c.Cache),
	}
}
