package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-looklook/app/order/api/internal/config"
	"go-zero-looklook/app/order/rpc/order"
	"go-zero-looklook/app/payment/rpc/payment"
	"go-zero-looklook/app/travel/rpc/homestayservice"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   order.Order
	PaymentRpc payment.Payment
	TravelRpc  homestayservice.HomestayService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		TravelRpc:  homestayservice.NewHomestayService(zrpc.MustNewClient(c.TravelRpcConf)),
		PaymentRpc: payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
	}
}
