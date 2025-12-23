package homestay

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/api/internal/convert"
	"go-zero-looklook/app/travel/rpc/homestayservice"

	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BusinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// boss all homestay room
func NewBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BusinessListLogic {
	return &BusinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BusinessListLogic) BusinessList(req *types.BusinessListReq) (resp *types.BusinessListResp, err error) {
	rpcResp, err := l.svcCtx.HomestayRpc.BusinessList(l.ctx, &homestayservice.BusinessListReq{
		LastId:             req.LastId,
		PageSize:           req.PageSize,
		HomestayBusinessId: req.HomestayBusinessId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "HomestayRpc.GuessList")
	}
	list := make([]types.Homestay, len(rpcResp.List))
	for i, item := range rpcResp.List {
		list[i] = convert.ConvertRpcHomestayToApiHomestay(item)
	}
	resp = &types.BusinessListResp{
		List: list,
	}

	return
}
