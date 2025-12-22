package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestayDetailLogic) HomestayDetail(in *pb.HomestayDetailReq) (*pb.HomestayDetailResp, error) {
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "Rpc.HomestayDetailLogic FindOne")
	}
	h := &pb.Homestay{}
	_ = copier.Copy(&h, homestay)

	return &pb.HomestayDetailResp{
		Homestay: h,
	}, nil
}
