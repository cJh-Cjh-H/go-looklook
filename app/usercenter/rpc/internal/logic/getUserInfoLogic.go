package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-looklook/app/usercenter/rpc/usercenter"
	"go-zero-looklook/pkg/xerr"

	"go-zero-looklook/app/usercenter/rpc/internal/svc"
	"go-zero-looklook/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	//查询user表
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GetUserInfo find user db err , id:%d , err:%v", in.Id, err)
	}
	var respUser usercenter.User
	_ = copier.Copy(&respUser, user)
	return &pb.GetUserInfoResp{
		User: &respUser,
	}, nil
}
