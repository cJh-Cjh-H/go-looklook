package homestayBusiness

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/model"
	"go-zero-looklook/app/travel/rpc/pb"
	"go-zero-looklook/app/usercenter/rpc/usercenter"
	"go-zero-looklook/pkg/xerr"

	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBusinessDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// boss detail
func NewHomestayBusinessDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBusinessDetailLogic {
	return &HomestayBusinessDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBusinessDetailLogic) HomestayBusinessDetail(req *types.HomestayBusinessDetailReq) (resp *types.HomestayBusinessDetailResp, err error) {
	homestayBusiness, err := l.svcCtx.HomestayRpc.HomestayBusinessDetail(l.ctx, &pb.HomestayBusinessDetailReq{
		Id: req.Id,
	})
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " HomestayBussinessDetail  FindOne db fail ,id  : %d , err : %v", req.Id, err)
	}
	if homestayBusiness == nil {
		return nil, errors.New("API.homestayBusiness == nil")
	}
	fmt.Printf("boos:%v\n", homestayBusiness.Boss)
	userResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
		Id: homestayBusiness.Boss.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Api.HomestayBusinessDetail.GetUserInfo")
	}
	user := userResp.User
	boss := types.HomestayBusinessBoss{
		Id:       homestayBusiness.Boss.Id,
		UserId:   homestayBusiness.Boss.UserId,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Rank:     homestayBusiness.Boss.Rank,
		Info:     user.Info,
	}
	return &types.HomestayBusinessDetailResp{
		Boss: boss,
	}, nil
}
