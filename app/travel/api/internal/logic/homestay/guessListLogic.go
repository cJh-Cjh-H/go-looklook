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

type GuessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// guess homestay room
func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuessListLogic {
	return &GuessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuessListLogic) GuessList(req *types.GuessListReq) (resp *types.GuessListResp, err error) {
	rpcResp, err := l.svcCtx.HomestayRpc.GuessList(l.ctx, &homestayservice.GuessListReq{})
	if err != nil {
		return nil, errors.Wrapf(err, "HomestayRpc.GuessList")
	}
	list := make([]types.Homestay, len(rpcResp.List))
	for i, item := range rpcResp.List {
		list[i] = convert.ConvertRpcHomestayToApiHomestay(item)
	}
	resp = &types.GuessListResp{
		List: list,
	}
	return
}
