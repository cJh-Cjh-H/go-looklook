package user

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-zero-looklook/app/usercenter/api/internal/svc"
	"go-zero-looklook/app/usercenter/api/internal/types"
	"go-zero-looklook/app/usercenter/model"
	"go-zero-looklook/app/usercenter/rpc/usercenter"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	fmt.Printf("reqp:%+v\n", req)
	registerResp, err := l.svcCtx.UsercenterRpc.Register(l.ctx, &usercenter.RegisterReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		AuthKey:  req.Mobile,
		AuthType: model.UserAuthTypeSystem,
	})
	fmt.Printf("registerResp:%+v\n", registerResp)
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	resp = &types.RegisterResp{
		AccessToken:  registerResp.AccessToken,
		AccessExpire: registerResp.AccessExpire,
		RefreshAfter: registerResp.RefreshAfter,
	}
	return resp, nil
}
