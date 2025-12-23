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
	"strings"
	"time"
)

var _ HomestayModel = (*customHomestayModel)(nil)

type (
	// HomestayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomestayModel.
	HomestayModel interface {
		homestayModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		TransInsert(ctx context.Context, session sqlx.Session, data *Homestay) (sql.Result, error)
		TransUpdate(ctx context.Context, session sqlx.Session, data *Homestay) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *Homestay) error
		SelectBuilder() squirrel.SelectBuilder
		DeleteSoft(ctx context.Context, session sqlx.Session, data *Homestay) error
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Homestay, error)
		FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Homestay, error)
		FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Homestay, int64, error)
		FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Homestay, error)
		FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Homestay, error)
		TransDelete(ctx context.Context, session sqlx.Session, id int64) error
		SelectBuilderWithJoin(asTableName string, joinType string, joinTable string, joinCondition string, whereCondition string, args ...interface{}) squirrel.SelectBuilder
		FindPageDIY(ctx context.Context, limit int64) ([]*Homestay, error)
	}

	customHomestayModel struct {
		*defaultHomestayModel
	}
)

func (m *defaultHomestayModel) FindPageDIY(ctx context.Context, limit int64) ([]*Homestay, error) {
	s := `SELECT
	h.*
		FROM homestay h
	LEFT JOIN (
		SELECT a.data_id, COUNT(*) AS activity_count
	FROM homestay_activity a
	WHERE a.row_status = 1 AND a.del_state=0
	GROUP BY a.data_id
	) ha ON h.id = ha.data_id
	LEFT JOIN (
		SELECT c.homestay_id,
		AVG(CAST(c.star->>'$.service' AS DECIMAL(10,2))) as avg
	FROM homestay_comment c
	WHERE  c.del_state = 0
	GROUP BY c.homestay_id
	) hc ON hc.homestay_id = h.id
	WHERE h.row_state = 1 AND h.del_state=0
	ORDER BY
	COALESCE(ha.activity_count, 0) DESC,
		COALESCE(CAST(hc.avg AS DECIMAL(10,2)), 0) DESC
	LIMIT ? ;
`
	var resp []*Homestay
	err := m.QueryRowsNoCache(&resp, s, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "Model.FindPageDIY.QueryRowNoCache")
	}
	return resp, err

}
func (m *defaultHomestayModel) SelectBuilderWithJoin(
	asTableName string,
	joinType string, // JOIN 类型：INNER、LEFT、RIGHT
	joinTable string, // 要 JOIN 的表名
	joinCondition string, // JOIN 条件
	whereCondition string, // WHERE 条件
	args ...interface{}, // 参数
) squirrel.SelectBuilder {

	baseQuery := squirrel.Select(fmt.Sprintf("%s.*", asTableName)).
		From(fmt.Sprintf("%s %s", m.table, asTableName))

	// 根据 JOIN 类型构建不同的 JOIN
	switch strings.ToUpper(joinType) {
	case "LEFT", "LEFT JOIN":
		baseQuery = baseQuery.LeftJoin(fmt.Sprintf("%s ON %s", joinTable, joinCondition))
	case "RIGHT", "RIGHT JOIN":
		baseQuery = baseQuery.RightJoin(fmt.Sprintf("%s ON %s", joinTable, joinCondition))
	default: // INNER JOIN
		baseQuery = baseQuery.Join(fmt.Sprintf("%s ON %s", joinTable, joinCondition))
	}

	if whereCondition != "" {
		baseQuery = baseQuery.Where(whereCondition, args...)
	}

	return baseQuery
}

func (m *defaultHomestayModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultHomestayModel) TransInsert(ctx context.Context, session sqlx.Session, data *Homestay) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	data.DelState = globalkey.DelStateNo
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, homestayRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Title, data.SubTitle, data.Banner, data.Info, data.PeopleNum, data.HomestayBusinessId, data.UserId, data.RowState, data.RowType, data.FoodInfo, data.FoodPrice, data.HomestayPrice, data.MarketHomestayPrice)
	}, looklookTravelHomestayIdKey)
}

func (m *defaultHomestayModel) TransUpdate(ctx context.Context, session sqlx.Session, data *Homestay) (sql.Result, error) {
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homestayRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Title, data.SubTitle, data.Banner, data.Info, data.PeopleNum, data.HomestayBusinessId, data.UserId, data.RowState, data.RowType, data.FoodInfo, data.FoodPrice, data.HomestayPrice, data.MarketHomestayPrice, data.Id)
	}, looklookTravelHomestayIdKey)
}

func (m *defaultHomestayModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *Homestay) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLookLookHomestayIdPrefix, data.Id)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, homestayRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Title, data.SubTitle, data.Banner, data.Info, data.PeopleNum, data.HomestayBusinessId, data.UserId, data.RowState, data.RowType, data.FoodInfo, data.FoodPrice, data.HomestayPrice, data.MarketHomestayPrice, data.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Title, data.SubTitle, data.Banner, data.Info, data.PeopleNum, data.HomestayBusinessId, data.UserId, data.RowState, data.RowType, data.FoodInfo, data.FoodPrice, data.HomestayPrice, data.MarketHomestayPrice, data.Id, oldVersion)
	}, looklookTravelHomestayIdKey)
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
func (m *defaultHomestayModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *Homestay) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(ctx, session, data); err != nil {
		return errors.Wrapf(errors.New("delete soft failed "), "HomestayModel delete err : %+v", err)
	}
	return nil
}

// FindSum 通用求和方法
func (m *defaultHomestayModel) FindSum(ctx context.Context, builder squirrel.SelectBuilder, field string) (float64, error) {

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
func (m *defaultHomestayModel) FindCount(ctx context.Context, builder squirrel.SelectBuilder, field string) (int64, error) {

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
func (m *defaultHomestayModel) FindAll(ctx context.Context, builder squirrel.SelectBuilder, orderBy string) ([]*Homestay, error) {

	builder = builder.Columns(homestayRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Homestay
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 分页查询的通用方法
func (m *defaultHomestayModel) FindPageListByPage(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Homestay, error) {

	//builder = builder.Columns(homestayRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	//query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	query, values, err := builder.Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Homestay
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPageWithTotal 分页查询的通用方法（带总数）
func (m *defaultHomestayModel) FindPageListByPageWithTotal(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Homestay, int64, error) {

	total, err := m.FindCount(ctx, builder, "id")
	if err != nil {
		return nil, 0, err
	}

	builder = builder.Columns(homestayRows)

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

	var resp []*Homestay
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		return nil, total, err
	}
}

// FindPageListByIdDESC 通过Id分页查询并降序排列
func (m *defaultHomestayModel) FindPageListByIdDESC(ctx context.Context, builder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Homestay, error) {

	builder = builder.Columns(homestayRows)

	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Homestay
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 通过Id分页查询并降序排列
func (m *defaultHomestayModel) FindPageListByIdASC(ctx context.Context, builder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Homestay, error) {

	builder = builder.Columns(homestayRows)

	if preMaxId > 0 {
		builder = builder.Where(" id > ? ", preMaxId)
	}

	query, values, err := builder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Homestay
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// SelectBuilder SQL 查询构建器的工厂方法
func (m *defaultHomestayModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}
func (m *defaultHomestayModel) TransDelete(ctx context.Context, session sqlx.Session, id int64) error {
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

// NewHomestayModel returns a model for the database table.
func NewHomestayModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) HomestayModel {
	return &customHomestayModel{
		defaultHomestayModel: newHomestayModel(conn, c, opts...),
	}
}
