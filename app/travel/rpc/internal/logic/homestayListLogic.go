package logic

import (
	"context"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿服务
func (l *HomestayListLogic) HomestayList(in *pb.HomestayListReq) (*pb.HomestayListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.HomestayListResp{}, nil
}
