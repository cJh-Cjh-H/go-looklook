package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-looklook/app/order/model"
	"go-zero-looklook/app/order/rpc/internal/config"
	"go-zero-looklook/app/travel/rpc/homestayservice"
)

type ServiceContext struct {
	Config config.Config

	TravelRpc          homestayservice.HomestayService
	RedisClient        *redis.Redis
	HomestayOrderModel model.HomestayOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		TravelRpc:          homestayservice.NewHomestayService(zrpc.MustNewClient(c.TravelRpcConf)),
		HomestayOrderModel: model.NewHomestayOrderModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		RedisClient:        redis.MustNewRedis(c.Redis.RedisConf),
	}
}
