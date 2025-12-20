package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-looklook/pkg/globalkey"
	"time"
)

var _ UserAuthModel = (*customUserAuthModel)(nil)

type (
	// UserAuthModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAuthModel.
	UserAuthModel interface {
		userAuthModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		TransInsert(ctx context.Context, session sqlx.Session, data *UserAuth) (sql.Result, error)
	}

	customUserAuthModel struct {
		*defaultUserAuthModel
	}
)

func (m *defaultUserAuthModel) TransInsert(ctx context.Context, session sqlx.Session, data *UserAuth) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	data.DelState = globalkey.DelStateNo
	lookLookUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLookLookUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	lookLookUserAuthIdKey := fmt.Sprintf("%s%v", cacheLookLookUserAuthIdPrefix, data.Id)
	lookLookUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLookLookUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, userAuthRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType)
	}, lookLookUserAuthAuthTypeAuthKeyKey, lookLookUserAuthIdKey, lookLookUserAuthUserIdAuthTypeKey)
	return ret, err
}

// Trans 把model层的事务方法暴露给logic层使用，因为logic层比model层更需要事务
func (m *defaultUserAuthModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})

}

// NewUserAuthModel returns a model for the database table.
func NewUserAuthModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserAuthModel {
	return &customUserAuthModel{
		defaultUserAuthModel: newUserAuthModel(conn, c, opts...),
	}
}
