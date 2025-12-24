package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/model"
	"go-zero-looklook/pkg/xerr"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodBossLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGoodBossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodBossLogic {
	return &GoodBossLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GoodBoss 民宿商家服务
func (l *GoodBossLogic) GoodBoss(in *pb.GoodBossReq) (*pb.GoodBossResp, error) {

	homestayActivityList, err := l.svcCtx.HomestayActivityModel.FindDiy(l.ctx, 10)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "get GoodBoss db err. rowType: %s ,err : %v", model.HomestayActivityGoodBusiType, err)
	}
	for _, t := range homestayActivityList {
		fmt.Printf("t:%v\n", t)
	}
	list := make([]*pb.HomestayBusinessBoss, len(homestayActivityList))
	for i, item := range homestayActivityList {
		list[i] = &pb.HomestayBusinessBoss{
			Id:       item.Id,
			UserId:   item.UserId,
			Nickname: "",
			Avatar:   "",
			Info:     "",
			Rank:     -1,
		}
	}
	return &pb.GoodBossResp{
		List: list,
	}, nil
}
