package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿评论服务
func (l *CommentListLogic) CommentList(in *pb.CommentListReq) (*pb.CommentListResp, error) {
	commentList, err := l.svcCtx.HomestayCommentModel.FindDIY(l.ctx, in.LastId, in.PageSize)
	if err != nil {
		return nil, errors.Wrapf(err, "Rpc.CommentListLogic FindDIY error.")
	}

	return &pb.CommentListResp{
		List: commentList,
	}, nil
}
