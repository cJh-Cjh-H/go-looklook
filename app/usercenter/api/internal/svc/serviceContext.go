package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-looklook/app/usercenter/api/internal/config"
	"go-zero-looklook/app/usercenter/rpc/usercenter"
)

type ServiceContext struct {
	Config        config.Config
	UsercenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
