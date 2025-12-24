package homestayBusiness

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/threading"
	"go-zero-looklook/app/travel/api/internal/svc"
	"go-zero-looklook/app/travel/api/internal/types"
	"go-zero-looklook/app/travel/rpc/pb"
	"go-zero-looklook/app/usercenter/rpc/usercenter"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodBossLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// good boss
func NewGoodBossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodBossLogic {
	return &GoodBossLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodBossLogic) GoodBoss(req *types.GoodBossReq) (resp *types.GoodBossResp, err error) {
	rpcResp, err := l.svcCtx.HomestayRpc.GoodBoss(l.ctx, &pb.GoodBossReq{})
	bosses := rpcResp.List
	if err != nil {
		return nil, errors.Wrapf(err, " API.GoodBoss req:%+v", req)
	}
	if len(bosses) == 0 {
		return nil, errors.Wrapf(err, " API.GoodBoss req:%+v not found")
	}
	//统一格式
	type result struct {
		index int
		user  *usercenter.User
		err   error
	}
	resultChan := make(chan result, len(bosses))

	//使用goSafe安全使用goroutine
	group := threading.NewRoutineGroup()
	for i := range bosses {
		index := i
		boss := bosses[index]
		group.RunSafe(func() {

			//调用userRpc
			userInfo, err2 := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
				Id: boss.UserId,
			})
			user := userInfo.User
			if err2 != nil {
				resultChan <- result{
					index: index,
					user:  nil,
					err:   errors.Wrapf(err2, " API.RunSafe.GetUserInfo"),
				}

				return
			}
			resultChan <- result{
				index: index,
				user:  user,
				err:   nil,
			}
		})

	}
	// 等待所有 goroutine 完成
	group.Wait()
	close(resultChan)
	for res := range resultChan {
		if res.err != nil {
			return nil, err
		}
		bosses[res.index].Nickname = res.user.Nickname
		bosses[res.index].Info = res.user.Info
		bosses[res.index].Avatar = res.user.Avatar
		bosses[res.index].UserId = res.user.Id
	}

	for _, v := range bosses {
		fmt.Printf("boss:%v\n", v)
	}
	list := make([]types.HomestayBusinessBoss, len(bosses))
	for i, boss := range bosses {
		list[i] = types.HomestayBusinessBoss{
			Id:       boss.Id,
			UserId:   boss.UserId,
			Nickname: boss.Nickname,
			Info:     boss.Info,
			Avatar:   boss.Avatar,
			Rank:     boss.Rank,
		}
	}

	resp = &types.GoodBossResp{
		List: list,
	}

	return resp, nil
}
