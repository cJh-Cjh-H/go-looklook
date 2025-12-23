package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/homestayservice"
	"go-zero-looklook/pkg/xerr"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BusinessListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BusinessListLogic {
	return &BusinessListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BusinessListLogic) BusinessList(in *pb.BusinessListReq) (*pb.BusinessListResp, error) {
	whereBuilder := l.svcCtx.HomestayModel.SelectBuilder().Where(squirrel.Eq{"homestay_business_id": in.HomestayBusinessId})
	homestays, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, whereBuilder, in.LastId, in.PageSize)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "HomestayBusinessId: %d ,err : %v", in.HomestayBusinessId, err)
	}
	list := make([]*homestayservice.Homestay, len(homestays))
	for i, homestay := range homestays {
		list[i] = &homestayservice.Homestay{
			Id:                  homestay.Id,
			Title:               homestay.Title,
			SubTitle:            homestay.SubTitle,
			Banner:              homestay.Banner,
			Info:                homestay.Info,
			PeopleNum:           homestay.PeopleNum,
			HomestayBusinessId:  homestay.HomestayBusinessId,
			UserId:              homestay.UserId,
			RowState:            homestay.RowState,
			RowType:             homestay.RowType,
			FoodInfo:            homestay.FoodInfo,
			FoodPrice:           float64(homestay.FoodPrice),
			HomestayPrice:       float64(homestay.HomestayPrice),
			MarketHomestayPrice: float64(homestay.MarketHomestayPrice),
		}
	}

	return &pb.BusinessListResp{
		List: list,
	}, nil
}
