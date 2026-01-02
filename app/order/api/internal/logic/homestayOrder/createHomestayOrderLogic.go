package homestayOrder

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-looklook/app/order/api/internal/svc"
	"go-zero-looklook/app/order/api/internal/types"
	"go-zero-looklook/app/order/rpc/order"
	"go-zero-looklook/app/travel/rpc/pb"
	"go-zero-looklook/pkg/ctxdata"
	"go-zero-looklook/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateHomestayOrderLogic 创建民宿订单
func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateHomestayOrderLogic) CreateHomestayOrder(req *types.CreateHomestayOrderReq) (resp *types.CreateHomestayOrderResp, err error) {
	// 1. 验证房源是否存在
	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &pb.HomestayDetailReq{
		Id: req.HomestayId,
	})
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, errors.Wrapf(err, "Api.TravelRpc.HomestayDetail")
	}
	if homestayResp.Homestay == nil || homestayResp.Homestay.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("API.homestay no exists"), "CreateHomestayOrder homestay no exists id : %d", req.HomestayId)
	}
	// 2. 从上下文中获取用户ID
	userId := ctxdata.GetUidFromCtx(l.ctx)
	// 3. 调用订单微服务的RPC接口创建订单
	rpcResp, err := l.svcCtx.OrderRpc.CreateHomestayOrder(l.ctx, &order.CreateHomestayOrderReq{
		HomestayId:    req.HomestayId,
		IsFood:        req.IsFood,
		LiveStartTime: req.LiveStartTime,
		LiveEndTime:   req.LiveEndTime,
		UserId:        userId,
		LivePeopleNum: req.LivePeopleNum,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "API.OrderRpc.CreateHomestayOrder")
	}
	// 4. 返回成功响应（只返回订单号）
	return &types.CreateHomestayOrderResp{
		OrderSn: rpcResp.Sn,
	}, nil
}
