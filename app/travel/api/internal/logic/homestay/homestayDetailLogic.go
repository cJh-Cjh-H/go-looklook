package homestay

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/homestayservice"

	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay room detail
func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayDetailLogic) HomestayDetail(req *types.HomestayDetailReq) (resp *types.HomestayDetailResp, err error) {
	homestayDetail, err := l.svcCtx.HomestayRpc.HomestayDetail(l.ctx, &homestayservice.HomestayDetailReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Api.homestayDetailLogic.HomestayDetail.HomestayDetail")
	}
	var h types.Homestay
	_ = copier.Copy(&h, homestayDetail.Homestay)
	resp = &types.HomestayDetailResp{
		Homestay: h,
	}
	return
}
