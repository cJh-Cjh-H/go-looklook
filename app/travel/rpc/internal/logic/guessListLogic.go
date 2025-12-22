package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/homestayservice"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GuessListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuessListLogic {
	return &GuessListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GuessListLogic) GuessList(in *pb.GuessListReq) (*pb.GuessListResp, error) {
	homestays, err := l.svcCtx.HomestayModel.FindPageDIY(l.ctx, 10)
	if err != nil {
		return nil, errors.Wrapf(err, "Rpc.GuessList.FindPageDIY")
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
	return &pb.GuessListResp{
		List: list,
	}, nil
}
