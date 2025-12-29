package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-zero-looklook/app/order/model"
	"go-zero-looklook/app/travel/rpc/homestayservice"
	"go-zero-looklook/pkg/redisLock"
	"go-zero-looklook/pkg/tool"
	"go-zero-looklook/pkg/uniqueid"
	"go-zero-looklook/pkg/xerr"
	"strings"
	"time"

	"go-zero-looklook/app/order/rpc/internal/svc"
	"go-zero-looklook/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿下订单
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {
	// ====== 1. 在这里首先加Redis分布式锁 ======
	lockKey := fmt.Sprintf("order_lock:homestay:%d:%d:%d",
		in.HomestayId, in.LiveStartTime, in.LiveEndTime)
	// 尝试获取分布式锁
	err := redisLock.Lock(l.svcCtx.RedisClient, lockKey, "1", 5)
	if err != nil {
		return nil, errors.Wrapf(err, "获取锁失败，请稍后重试")
	}
	
	//1、Create Order
	if in.LiveEndTime <= in.LiveStartTime {
		return nil, errors.Wrapf(xerr.NewErrMsg("Stay at least one night"), "Place an order at a B&B. The end time of your stay must be greater than the start time. in : %+v", in)
	}

	resp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &homestayservice.HomestayDetailReq{
		Id: in.HomestayId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to query the record"), "Failed to query the record  rpc HomestayDetail fail , homestayId : %d , err : %v", in.HomestayId, err)
	}
	if resp.Homestay == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("This record does not exist"), "This record does not exist , homestayId : %d ", in.HomestayId)
	}

	var cover string //Get the cover...
	if len(resp.Homestay.Banner) > 0 {
		cover = strings.Split(resp.Homestay.Banner, ",")[0]
	}

	order := new(model.HomestayOrder)
	order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_HOMESTAY_ORDER)
	order.UserId = in.UserId
	order.HomestayId = in.HomestayId
	order.Title = resp.Homestay.Title
	order.SubTitle = resp.Homestay.SubTitle
	order.Cover = cover
	order.Info = resp.Homestay.Info
	order.PeopleNum = resp.Homestay.PeopleNum
	order.RowType = resp.Homestay.RowType
	order.HomestayPrice = int64(resp.Homestay.HomestayPrice)
	order.MarketHomestayPrice = int64(resp.Homestay.MarketHomestayPrice)
	order.HomestayBusinessId = resp.Homestay.HomestayBusinessId
	order.HomestayUserId = resp.Homestay.UserId
	order.LivePeopleNum = in.LivePeopleNum
	order.TradeState = model.HomestayOrderTradeStateWaitPay
	order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	order.Remark = in.Remark
	order.FoodInfo = resp.Homestay.FoodInfo
	order.FoodPrice = int64(resp.Homestay.FoodPrice)
	order.LiveStartDate = time.Unix(in.LiveStartTime, 0)
	order.LiveEndDate = time.Unix(in.LiveEndTime, 0)

	liveDays := int64(order.LiveEndDate.Sub(order.LiveStartDate).Seconds() / 86400) //Stayed a few days in total

	order.HomestayTotalPrice = int64(int64(resp.Homestay.HomestayPrice) * liveDays) //Calculate the total price of the B&B
	if in.IsFood {
		order.NeedFood = model.HomestayOrderNeedFoodYes
		//Calculate the total price of the meal.
		order.FoodTotalPrice = int64(int64(resp.Homestay.FoodPrice) * in.LivePeopleNum * liveDays)
	}

	order.OrderTotalPrice = order.HomestayTotalPrice + order.FoodTotalPrice //Calculate total order price.

	_, err = l.svcCtx.HomestayOrderModel.Insert(l.ctx, nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Order Database Exception order : %+v , err: %v", order, err)
	}
	return &pb.CreateHomestayOrderResp{
		Sn: order.Sn,
	}, nil
	//2、Delayed closing of order tasks.
	//payload, err := json.Marshal(jobtype.DeferCloseHomestayOrderPayload{Sn: order.Sn})
	//if err != nil {
	//	logx.WithContext(l.ctx).Errorf("create defer close order task json Marshal fail err :%+v , sn : %s", err, order.Sn)
	//} else {
	//	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.DeferCloseHomestayOrder, payload), asynq.ProcessIn(CloseOrderTimeMinutes*time.Minute))
	//	if err != nil {
	//		logx.WithContext(l.ctx).Errorf("create defer close order task insert queue fail err :%+v , sn : %s", err, order.Sn)
	//	}
	//}

}
