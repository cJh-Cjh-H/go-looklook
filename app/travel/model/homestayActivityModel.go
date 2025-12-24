package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-looklook/app/usercenter/model"
	"go-zero-looklook/pkg/globalkey"
	"time"
)

var _ HomestayActivityModel = (*customHomestayActivityModel)(nil)

type (
	// HomestayActivityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomestayActivityModel.
	HomestayActivityModel interface {
		homestayActivityModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		TransInsert(ctx context.Context, session sqlx.Session, data *HomestayActivity) (sql.Result, error)
		TransUpdate(ctx context.Context, session sqlx.Session, data *HomestayActivity) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *HomestayActivity) error
		SelectBuilder() squirrel.SelectBuilder
		DeleteSoft(ctx context.Context, session sqlx.Session, data *HomestayActivity) error
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*HomestayActivity, error)
		FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayActivity, error)
		FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayActivity, int64, error)
		FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*HomestayActivity, error)
		FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*HomestayActivity, error)
		TransDelete(ctx context.Context, session sqlx.Session, id int64) error
		FindDiy(ctx context.Context, l int) ([]*HomestayBusinessBoss, error)
	}

	customHomestayActivityModel struct {
		*defaultHomestayActivityModel
	}
)

func (m *defaultHomestayActivityModel) FindDiy(ctx context.Context, l int) ([]*HomestayBusinessBoss, error) {
	s := `SELECT
    ha.id,s.user_id
FROM homestay_activity ha
         LEFT JOIN (
    SELECT h.user_id,h.id
    FROM homestay h
    WHERE h.row_state = 1 AND h.del_state=0
) as s  ON s.id = ha.data_id
WHERE ha.row_status = 1 AND ha.del_state=0 AND ha.row_type='goodBusiness'
LIMIT ? ;`
	var resp []*HomestayBusinessBoss

	err := m.QueryRowsNoCache(&resp, s, l)
	if err != nil {
		return nil, errors.Wrapf(err, "Model.defaultHomestayActivityModel.FindDiy.QueryRowNoCache")
	}
	return resp, err
}
func (m *defaultHomestayActivityModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultHomestayActivityModel) TransInsert(ctx context.Context, session sqlx.Session, data *HomestayActivity) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	data.DelState = globalkey.DelStateNo
	looklookTravelHomestayActivityIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayActivityIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, homestayActivityRowsExpectAutoSet)

		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.RowType, data.DataId, data.RowStatus, data.Version)
	}, looklookTravelHomestayActivityIdKey)
}

func (m *defaultHomestayActivityModel) TransUpdate(ctx context.Context, session sqlx.Session, data *HomestayActivity) (sql.Result, error) {
	looklookTravelHomestayActivityIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayActivityIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homestayActivityRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.RowType, data.DataId, data.RowStatus, data.Version, data.Id)
	}, looklookTravelHomestayActivityIdKey)
}

func (m *defaultHomestayActivityModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *HomestayActivity) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	looklookTravelHomestayActivityIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayActivityIdPrefix, data.Id)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, homestayActivityRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.RowType, data.DataId, data.RowStatus, data.Version, data.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.RowType, data.DataId, data.RowStatus, data.Version, data.Id, oldVersion)
	}, looklookTravelHomestayActivityIdKey)
	if err != nil {
		return err
	}
	updateCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if updateCount == 0 {
		return model.ErrNoRowsUpdate
	}

	return nil
}

// DeleteSoft 软删除
func (m *defaultHomestayActivityModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *HomestayActivity) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(ctx, session, data); err != nil {
		return errors.Wrapf(errors.New("delete soft failed "), "HomestayModel delete err : %+v", err)
	}
	return nil
}

// FindSum 通用求和方法
func (m *defaultHomestayActivityModel) FindSum(ctx context.Context, builder squirrel.SelectBuilder, field string) (float64, error) {

	if len(field) == 0 {
		return 0, errors.Wrapf(errors.New("FindSum Least One Field"), "FindSum Least One Field")
	}

	builder = builder.Columns("IFNULL(SUM(" + field + "),0)")

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindCount 通用的查询数量方法
func (m *defaultHomestayActivityModel) FindCount(ctx context.Context, builder squirrel.SelectBuilder, field string) (int64, error) {

	if len(field) == 0 {
		return 0, errors.Wrapf(errors.New("FindCount Least One Field"), "FindCount Least One Field")
	}

	builder = builder.Columns("COUNT(" + field + ")")

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindAll 通用的查询所有记录方法
func (m *defaultHomestayActivityModel) FindAll(ctx context.Context, builder squirrel.SelectBuilder, orderBy string) ([]*HomestayActivity, error) {

	builder = builder.Columns(homestayActivityRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayActivity
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 分页查询的通用方法
func (m *defaultHomestayActivityModel) FindPageListByPage(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayActivity, error) {

	builder = builder.Columns(homestayActivityRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayActivity
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPageWithTotal 分页查询的通用方法（带总数）
func (m *defaultHomestayActivityModel) FindPageListByPageWithTotal(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayActivity, int64, error) {

	total, err := m.FindCount(ctx, builder, "id")
	if err != nil {
		return nil, 0, err
	}

	builder = builder.Columns(homestayActivityRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, total, err
	}

	var resp []*HomestayActivity
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		return nil, total, err
	}
}

// FindPageListByIdDESC 通过Id分页查询并降序排列
func (m *defaultHomestayActivityModel) FindPageListByIdDESC(ctx context.Context, builder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*HomestayActivity, error) {

	builder = builder.Columns(homestayActivityRows)

	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayActivity
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 通过Id分页查询并降序排列
func (m *defaultHomestayActivityModel) FindPageListByIdASC(ctx context.Context, builder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*HomestayActivity, error) {

	builder = builder.Columns(homestayActivityRows)

	if preMaxId > 0 {
		builder = builder.Where(" id > ? ", preMaxId)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayActivity
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// SelectBuilder SQL 查询构建器的工厂方法
func (m *defaultHomestayActivityModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}
func (m *defaultHomestayActivityModel) TransDelete(ctx context.Context, session sqlx.Session, id int64) error {
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, looklookTravelHomestayIdKey)
	return err
}

// NewHomestayActivityModel returns a model for the database table.
func NewHomestayActivityModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) HomestayActivityModel {
	return &customHomestayActivityModel{
		defaultHomestayActivityModel: newHomestayActivityModel(conn, c, opts...),
	}
}
