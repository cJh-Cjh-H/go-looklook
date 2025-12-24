package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBusinessDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayBusinessDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBusinessDetailLogic {
	return &HomestayBusinessDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestayBusinessDetailLogic) HomestayBusinessDetail(in *pb.HomestayBusinessDetailReq) (*pb.HomestayBusinessDetailResp, error) {
	business, err := l.svcCtx.HomestayBusinessModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "Rpc.HomestayBusinessModel.FindOne(%d)", in.Id)
	}
	boss := &pb.HomestayBusinessBoss{
		Id:       business.Id,
		UserId:   business.UserId,
		Nickname: "",
		Avatar:   "",
		Info:     business.Info,
		Rank:     -1,
	}
	return &pb.HomestayBusinessDetailResp{
		Boss: boss,
	}, nil
}
