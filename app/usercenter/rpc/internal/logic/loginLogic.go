package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/usercenter/rpc/usercenter"

	"go-zero-looklook/app/usercenter/rpc/internal/svc"
	"go-zero-looklook/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	//查找用户是否存在
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.AuthKey)
	if err != nil {
		return nil, err
	}
	//比较密码
	if user.Password != in.Password {
		return nil, ErrUsernamePwdError
	}

	//生成token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: user.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", user.Id)
	}

	return &pb.LoginResp{
		AccessExpire: tokenResp.AccessExpire,
		AccessToken:  tokenResp.AccessToken,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
