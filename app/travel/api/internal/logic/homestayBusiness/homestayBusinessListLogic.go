package homestayBusiness

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/pb"

	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBusinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// business list
func NewHomestayBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBusinessListLogic {
	return &HomestayBusinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBusinessListLogic) HomestayBusinessList(req *types.HomestayBusinessListReq) (resp *types.HomestayBusinessListResp, err error) {
	rpcResp, err := l.svcCtx.HomestayRpc.HomestayBusinessList(l.ctx, &pb.HomestayBusinessListReq{
		LastId:   req.LastId,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Api.HomestayBusinessListLogic.HomestayBusinessList")
	}
	if len(rpcResp.List) == 0 {
		return nil, errors.Wrapf(err, "Api.len(rpcResp.List) == 0.HomestayBusinessListLogic.HomestayBusinessList")
	}
	list := make([]types.HomestayBusinessListInfo, len(rpcResp.List))
	for i, item := range rpcResp.List {
		res := item.HomestayBusiness
		hb := types.HomestayBusiness{
			Id:        res.Id,
			Title:     res.Title,
			Info:      res.Info,
			Tags:      res.Tags,
			Cover:     res.Cover,
			Star:      res.Star,
			IsFav:     res.IsFav,
			HeaderImg: res.HeaderImg,
		}

		list[i] = types.HomestayBusinessListInfo{
			HomestayBusiness: hb,
			SellMonth:        item.SellMonth,
			PersonConsume:    item.PersonConsume,
		}
	}
	resp = &types.HomestayBusinessListResp{
		List: list,
	}
	return resp, nil
}
