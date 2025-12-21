package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"go-zero-looklook/app/travel/model"

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
	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityPreferredType,
		"row_status": model.HomestayActivityUpStatus,
	})
	homestays, err := l.svcCtx.HomestayModel.FindByActivity(
		l.ctx,
		model.HomestayActivityPreferredType,
		model.HomestayActivityUpStatus,
		in.Page,
		in.PageSize,
	)

	return &pb.HomestayListResp{}, nil
}
