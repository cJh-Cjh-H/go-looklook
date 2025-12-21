package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-looklook/app/travel/api/internal/config"
	"go-zero-looklook/app/travel/rpc/homestayservice"
	_ "go-zero-looklook/app/travel/rpc/homestayservice"
	"go-zero-looklook/app/usercenter/rpc/usercenter"
)

type ServiceContext struct {
	Config        config.Config
	HomestayRpc   homestayservice.HomestayService
	UsercenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		HomestayRpc:   homestayservice.NewHomestayService(zrpc.MustNewClient(c.TravelRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
