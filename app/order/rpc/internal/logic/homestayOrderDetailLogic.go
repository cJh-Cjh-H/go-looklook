package logic

import (
	"context"

	"go-zero-looklook/app/order/rpc/internal/svc"
	"go-zero-looklook/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayOrderDetailLogic {
	return &HomestayOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿订单详情
func (l *HomestayOrderDetailLogic) HomestayOrderDetail(in *pb.HomestayOrderDetailReq) (*pb.HomestayOrderDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.HomestayOrderDetailResp{}, nil
}
