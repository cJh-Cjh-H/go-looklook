package homestayOrder

import (
	"context"
	"go-zero-looklook/app/order/rpc/order"
	"go-zero-looklook/pkg/ctxdata"

	"go-zero-looklook/app/order/api/internal/svc"
	"go-zero-looklook/app/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserHomestayOrderListLogic 用户订单列表
func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderListLogic) UserHomestayOrderList(req *types.UserHomestayOrderListReq) (resp *types.UserHomestayOrderListResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	rpcResp, err := l.svcCtx.OrderRpc.UserHomestayOrderList(l.ctx, &order.UserHomestayOrderListReq{
		UserId:      userId,
		TraderState: req.TradeState,
		PageSize:    req.PageSize,
		LastId:      req.LastId,
	})
	list := make([]types.UserHomestayOrderListView, len(rpcResp.List))
	for i, o := range rpcResp.List {
		list[i] = types.UserHomestayOrderListView{
			Sn:              o.Sn,
			Title:           o.Title,
			SubTitle:        o.SubTitle,
			HomestayId:      o.HomestayId,
			Cover:           o.Cover,
			OrderTotalPrice: float64(o.OrderTotalPrice),
			CreateTime:      o.CreateTime,
			TradeState:      o.TradeState,
			LiveStartDate:   o.LiveStartDate,
			LiveEndDate:     o.LiveEndDate,
			TradeCode:       o.TradeCode,
		}
	}
	return &types.UserHomestayOrderListResp{
		List: list,
	}, nil

}
