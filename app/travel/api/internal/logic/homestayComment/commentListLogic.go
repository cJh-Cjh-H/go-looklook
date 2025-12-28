package homestayComment

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/rpc/pb"

	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay comment list
func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	rpcResp, err := l.svcCtx.HomestayRpc.CommentList(l.ctx, &pb.CommentListReq{
		LastId:   req.LastId,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "API.CommentListLogic.CommentList")
	}
	list := make([]types.HomestayComment, len(rpcResp.List))
	for i, comment := range rpcResp.List {
		list[i] = types.HomestayComment{
			Id:         comment.Id,
			UserId:     comment.UserId,
			Star:       comment.Star,
			Content:    comment.Content,
			HomestayId: comment.HomestayId,
			Avatar:     comment.Avatar,
			Nickname:   comment.Nickname,
		}
	}
	resp = &types.CommentListResp{
		List: list,
	}
	return
}
