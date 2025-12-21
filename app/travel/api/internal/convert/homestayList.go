package convert

import (
	"go-zero-looklook/app/travel/api/internal/types"
	"go-zero-looklook/app/travel/rpc/pb"
)

// ConvertRpcHomestayToApiHomestay internal/logic/homestay/homestaylistlogic.go
func ConvertRpcHomestayToApiHomestay(rpcHomestay *pb.Homestay) types.Homestay {
	return types.Homestay{
		Id:                  rpcHomestay.Id,
		Title:               rpcHomestay.Title,
		SubTitle:            rpcHomestay.SubTitle,
		Banner:              rpcHomestay.Banner,
		Info:                rpcHomestay.Info,
		PeopleNum:           rpcHomestay.PeopleNum,
		HomestayBusinessId:  rpcHomestay.HomestayBusinessId,
		UserId:              rpcHomestay.UserId,
		RowState:            rpcHomestay.RowState,
		RowType:             rpcHomestay.RowType,
		FoodInfo:            rpcHomestay.FoodInfo,
		FoodPrice:           rpcHomestay.FoodPrice,
		HomestayPrice:       rpcHomestay.HomestayPrice,
		MarketHomestayPrice: rpcHomestay.MarketHomestayPrice,
	}
}
