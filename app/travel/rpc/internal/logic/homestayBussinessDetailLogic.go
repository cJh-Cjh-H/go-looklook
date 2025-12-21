package logic

import (
	"context"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayBussinessDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessDetailLogic {
	return &HomestayBussinessDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestayBussinessDetailLogic) HomestayBussinessDetail(in *pb.HomestayBussinessDetailReq) (*pb.HomestayBussinessDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.HomestayBussinessDetailResp{}, nil
}
