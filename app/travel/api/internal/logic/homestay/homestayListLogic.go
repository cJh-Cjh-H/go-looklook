package homestay

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/api/internal/convert"
	"go-zero-looklook/app/travel/rpc/homestayservice"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"
	_ "go-zero-looklook/app/travel/rpc/homestayservice"
)

type HomestayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay room list
func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayListLogic) HomestayList(req *types.HomestayListReq) (resp *types.HomestayListResp, err error) {
	rpcResp, err := l.svcCtx.HomestayRpc.HomestayList(l.ctx, &homestayservice.HomestayListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "HomestayListLogic.HomestayList")
	}
	list := make([]types.Homestay, len(rpcResp.List))
	for i, item := range rpcResp.List {
		list[i] = convert.ConvertRpcHomestayToApiHomestay(item)
	}
	resp = &types.HomestayListResp{
		List: list,
	}
	return resp, nil
}
