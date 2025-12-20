package user

import (
	"context"
	"go-zero-looklook/app/usercenter/model"
	"go-zero-looklook/app/usercenter/rpc/usercenter"

	"go-zero-looklook/app/usercenter/api/internal/svc"
	"go-zero-looklook/app/usercenter/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResp, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &usercenter.LoginReq{
		Password: req.Password,
		AuthType: model.UserAuthTypeSystem,
		AuthKey:  req.Mobile,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.LoginResp{
		AccessToken:  loginResp.AccessToken,
		AccessExpire: loginResp.AccessExpire,
		RefreshAfter: loginResp.RefreshAfter,
	}
	return resp, nil
}
