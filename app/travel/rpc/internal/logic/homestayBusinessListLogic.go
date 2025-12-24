package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/pkg/xerr"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBusinessListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBusinessListLogic {
	return &HomestayBusinessListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestayBusinessListLogic) HomestayBusinessList(in *pb.HomestayBusinessListReq) (*pb.HomestayBusinessListResp, error) {
	whereBuilder := l.svcCtx.HomestayBusinessModel.SelectBuilder()
	businessPage, err := l.svcCtx.HomestayBusinessModel.FindPageListByIdDESC(l.ctx, whereBuilder, in.LastId, in.PageSize)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "HomestayBussinessList FindPageListByIdDESC db fail ,  req : %+v , err:%v", in, err)
	}
	if len(businessPage) == 0 {
		return nil, errors.Wrapf(errors.New("businessPage == nil"), "businessPage == nil")
	}
	// 转换响应
	list := make([]*pb.HomestayBusinessListInfo, len(businessPage))
	for i, business := range businessPage {
		list[i] = &pb.HomestayBusinessListInfo{}
		res := &pb.HomestayBusiness{
			Id:        business.Id,
			Title:     business.Title,
			Info:      business.Info,
			Tags:      business.Tags,
			Cover:     business.Cover,
			Star:      business.Star,
			IsFav:     0,
			HeaderImg: business.HeaderImg,
		}

		list[i].HomestayBusiness = res
		list[i].SellMonth = 0
		list[i].PersonConsume = 0
	}

	return &pb.HomestayBusinessListResp{
		List: list,
	}, nil
}
