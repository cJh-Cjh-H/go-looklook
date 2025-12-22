package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-looklook/app/travel/model"
	"go-zero-looklook/app/travel/rpc/homestayservice"
	"go-zero-looklook/pkg/globalkey"

	"go-zero-looklook/app/travel/rpc/internal/svc"
	"go-zero-looklook/app/travel/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿服务
func (l *HomestayListLogic) HomestayList(in *pb.HomestayListReq) (*pb.HomestayListResp, error) {
	//获取join查询的builder
	builderJoin := l.svcCtx.HomestayModel.SelectBuilderWithJoin(
		"h",
		"INNER",
		"homestay_activity a",
		"a.data_id = h.id",
		"a.row_type = ? AND a.row_status = ? AND a.del_state = ? AND h.del_state = ?",
		model.HomestayActivitySeasonType,
		model.HomestayActivityUpStatus,
		globalkey.DelStateNo,
		globalkey.DelStateNo,
	)
	//利用join的builder来进行分页查询
	homestaysPage, err := l.svcCtx.HomestayModel.FindPageListByPage(l.ctx, builderJoin, in.Page, in.PageSize, "")
	if err != nil {
		return nil, errors.Wrapf(err, "Rpc.homestayListLogic.HomestayList.FindPageListByPage")
	}
	// 转换响应
	list := make([]*homestayservice.Homestay, len(homestaysPage))
	for i, homestay := range homestaysPage {
		list[i] = &homestayservice.Homestay{
			Id:                  homestay.Id,
			Title:               homestay.Title,
			SubTitle:            homestay.SubTitle,
			Banner:              homestay.Banner,
			Info:                homestay.Info,
			PeopleNum:           homestay.PeopleNum,
			HomestayBusinessId:  homestay.HomestayBusinessId,
			UserId:              homestay.UserId,
			RowState:            homestay.RowState,
			RowType:             homestay.RowType,
			FoodInfo:            homestay.FoodInfo,
			FoodPrice:           float64(homestay.FoodPrice),
			HomestayPrice:       float64(homestay.HomestayPrice),
			MarketHomestayPrice: float64(homestay.MarketHomestayPrice),
		}
	}

	return &pb.HomestayListResp{
		List: list,
	}, nil
}
