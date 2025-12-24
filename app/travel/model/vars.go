package model

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrNoRowsUpdate = errors.New("update db no rows change")

// 民宿活动类型
var HomestayActivitySeasonType = "season_discount"      //季节民宿
var HomestayActivityPreferredType = "preferredHomestay" //优选民宿
var HomestayActivityGoodBusiType = "goodBusiness"       //最佳房东

// 民宿活动上下架

var HomestayActivityDownStatus int64 = 0 //下架
var HomestayActivityUpStatus int64 = 1   //上架
type HomestayBusinessBoss struct {
	Id     int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId int64 `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
}
