package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"go-zero-looklook/app/order/model"
	"go-zero-looklook/pkg/xerr"

	"go-zero-looklook/app/order/rpc/internal/svc"
	"go-zero-looklook/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户民宿订单
func (l *UserHomestayOrderListLogic) UserHomestayOrderList(in *pb.UserHomestayOrderListReq) (*pb.UserHomestayOrderListResp, error) {
	whereBuilder := l.svcCtx.HomestayOrderModel.SelectBuilder().Where(squirrel.Eq{"user_id": in.UserId})
	//There are supported states in the filter, otherwise return all
	if in.TraderState >= model.HomestayOrderTradeStateCancel && in.TraderState <= model.HomestayOrderTradeStateExpire {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"trade_state": in.TraderState})
	}

	resp, err := l.svcCtx.HomestayOrderModel.FindPageListByIdDESC(l.ctx, whereBuilder, in.LastId, in.PageSize)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Failed to get user's homestay order err : %v , in :%+v", err, in)
	}
	list := make([]*pb.HomestayOrder, len(resp))
	for i, o := range resp {
		list[i] = &pb.HomestayOrder{
			Id:                  o.Id,
			Sn:                  o.Sn,
			UserId:              o.UserId,
			HomestayId:          o.HomestayId,
			Title:               o.Title,
			SubTitle:            o.SubTitle,
			Cover:               o.Cover,
			Info:                o.Info,
			PeopleNum:           o.PeopleNum,
			RowType:             o.RowType,
			FoodInfo:            o.FoodInfo,
			FoodPrice:           o.FoodPrice,
			HomestayPrice:       o.HomestayPrice,
			MarketHomestayPrice: o.MarketHomestayPrice,
			HomestayBusinessId:  o.HomestayBusinessId,
			HomestayUserId:      o.HomestayUserId,
			LiveStartDate:       o.LiveStartDate.Unix(),
			LiveEndDate:         o.LiveEndDate.Unix(),
			LivePeopleNum:       o.LivePeopleNum,
			TradeState:          o.TradeState,
			TradeCode:           o.TradeCode,
			Remark:              o.Remark,
			OrderTotalPrice:     o.OrderTotalPrice,
			FoodTotalPrice:      o.FoodTotalPrice,
			HomestayTotalPrice:  o.HomestayTotalPrice,
			CreateTime:          o.CreateTime.Unix(),
			NeedFood:            o.NeedFood,
		}
	}
	return &pb.UserHomestayOrderListResp{
		List: list,
	}, nil
}
