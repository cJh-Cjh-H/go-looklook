package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-looklook/app/travel/rpc/pb"
	"go-zero-looklook/app/usercenter/model"
	"go-zero-looklook/pkg/globalkey"
	"time"
)

var _ HomestayCommentModel = (*customHomestayCommentModel)(nil)

type (
	// HomestayCommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomestayCommentModel.
	HomestayCommentModel interface {
		homestayCommentModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		TransInsert(ctx context.Context, session sqlx.Session, data *HomestayComment) (sql.Result, error)
		TransUpdate(ctx context.Context, session sqlx.Session, data *HomestayComment) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *HomestayComment) error
		SelectBuilder() squirrel.SelectBuilder
		DeleteSoft(ctx context.Context, session sqlx.Session, data *HomestayComment) error
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*HomestayComment, error)
		FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayComment, error)
		FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayComment, int64, error)
		FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*HomestayComment, error)
		FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*HomestayComment, error)
		TransDelete(ctx context.Context, session sqlx.Session, id int64) error
		FindDIY(ctx context.Context, lastId int64, pageSize int64) ([]*pb.HomestayComment, error)
	}

	customHomestayCommentModel struct {
		*defaultHomestayCommentModel
	}
)

func (m *defaultHomestayCommentModel) FindDIY(ctx context.Context, lastId int64, pageSize int64) ([]*pb.HomestayComment, error) {
	s := `SELECT c.id,
        c.homestay_id,
        c.content,  -- 放到第3位，对应pb.HomestayComment.Content
        AVG(
            (CAST(c.star->>'$.view' AS DECIMAL(3,2)) +
             CAST(c.star->>'$.clean' AS DECIMAL(3,2)) +
             CAST(c.star->>'$.service' AS DECIMAL(3,2)) +
             CAST(c.star->>'$.location' AS DECIMAL(3,2))
            ) / 4.0
        ) as star,
        c.user_id,
        u.nickname,
        u.avatar
FROM homestay_comment c
    INNER JOIN user u ON c.user_id = u.id
where c.del_state=0 and u.del_state=0 and c.id>?
GROUP BY c.id, c.homestay_id, c.user_id, c.content, u.nickname, u.avatar
limit ?;
`
	var resp []*pb.HomestayComment
	err := m.QueryRowsNoCache(&resp, s, lastId, pageSize)
	if err != nil {
		return nil, errors.Wrapf(err, "Model.FindPageDIY.QueryRowNoCache")
	}
	return resp, err
}
func (m *defaultHomestayCommentModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultHomestayCommentModel) TransInsert(ctx context.Context, session sqlx.Session, data *HomestayComment) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	data.DelState = globalkey.DelStateNo
	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayCommentIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, homestayCommentRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.HomestayId, data.UserId, data.Content, data.Star, data.Version)
	}, looklookTravelHomestayCommentIdKey)
}

func (m *defaultHomestayCommentModel) TransUpdate(ctx context.Context, session sqlx.Session, data *HomestayComment) (sql.Result, error) {
	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayCommentIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homestayCommentRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.HomestayId, data.UserId, data.Content, data.Star, data.Version, data.Id)
	}, looklookTravelHomestayCommentIdKey)
}

func (m *defaultHomestayCommentModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *HomestayComment) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayCommentIdPrefix, data.Id)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, homestayCommentRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.HomestayId, data.UserId, data.Content, data.Star, data.Version, data.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.HomestayId, data.UserId, data.Content, data.Star, data.Version, data.Id, oldVersion)
	}, looklookTravelHomestayCommentIdKey)
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
func (m *defaultHomestayCommentModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *HomestayComment) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(ctx, session, data); err != nil {
		return errors.Wrapf(errors.New("delete soft failed "), "HomestayModel delete err : %+v", err)
	}
	return nil
}

// FindSum 通用求和方法
func (m *defaultHomestayCommentModel) FindSum(ctx context.Context, builder squirrel.SelectBuilder, field string) (float64, error) {

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
func (m *defaultHomestayCommentModel) FindCount(ctx context.Context, builder squirrel.SelectBuilder, field string) (int64, error) {

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
func (m *defaultHomestayCommentModel) FindAll(ctx context.Context, builder squirrel.SelectBuilder, orderBy string) ([]*HomestayComment, error) {

	builder = builder.Columns(homestayCommentRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayComment
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 分页查询的通用方法
func (m *defaultHomestayCommentModel) FindPageListByPage(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayComment, error) {

	builder = builder.Columns(homestayCommentRows)

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

	var resp []*HomestayComment
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPageWithTotal 分页查询的通用方法（带总数）
func (m *defaultHomestayCommentModel) FindPageListByPageWithTotal(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*HomestayComment, int64, error) {

	total, err := m.FindCount(ctx, builder, "id")
	if err != nil {
		return nil, 0, err
	}

	builder = builder.Columns(homestayCommentRows)

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

	var resp []*HomestayComment
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		return nil, total, err
	}
}

// FindPageListByIdDESC 通过Id分页查询并降序排列
func (m *defaultHomestayCommentModel) FindPageListByIdDESC(ctx context.Context, builder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*HomestayComment, error) {

	builder = builder.Columns(homestayCommentRows)

	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayComment
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 通过Id分页查询并降序排列
func (m *defaultHomestayCommentModel) FindPageListByIdASC(ctx context.Context, builder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*HomestayComment, error) {

	builder = builder.Columns(homestayCommentRows)

	if preMaxId > 0 {
		builder = builder.Where(" id > ? ", preMaxId)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*HomestayComment
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// SelectBuilder SQL 查询构建器的工厂方法
func (m *defaultHomestayCommentModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}
func (m *defaultHomestayCommentModel) TransDelete(ctx context.Context, session sqlx.Session, id int64) error {
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

// NewHomestayCommentModel returns a model for the database table.
func NewHomestayCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) HomestayCommentModel {
	return &customHomestayCommentModel{
		defaultHomestayCommentModel: newHomestayCommentModel(conn, c, opts...),
	}
}
