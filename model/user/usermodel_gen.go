// Code generated by goctl. DO NOT EDIT!

package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userFieldNames          = builder.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUserIdPrefix       = "cache:user:id:"
	cacheUserAddressPrefix  = "cache:user:address:"
	cacheUserTwitterPrefix  = "cache:user:twitter:"
	cacheUserDiscordPrefix  = "cache:user:discord:"
	cacheUserUsernamePrefix = "cache:user:username:"
	cacheUserNNSPrefix      = "cache:user:nns:"
	cacheUserQueryPrefix    = "cache:user:genealquery:"
)

type (
	userModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*User, error)
		FindOneByAddress(ctx context.Context, address string) (*User, error)
		FindBySQL(ctx context.Context, sql string) (*[]User, error)
		FindOneByTwitter(ctx context.Context, twitter string) (*User, error)
		FindOneByDiscord(ctx context.Context, discord string) (*User, error)
		FindOneByUserName(ctx context.Context, UserName string) (*User, error)
		FindOneByNNS(ctx context.Context, NNS string) (*User, error)
		Update(ctx context.Context, newData *User) error
		Delete(ctx context.Context, id int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id            int64          `db:"id"`
		Username      sql.NullString `db:"username"`
		Bio           sql.NullString `db:"bio"`
		Address       string         `db:"address"`
		NNS           sql.NullString `db:"nns"`
		Email         sql.NullString `db:"email"`
		Twitter       sql.NullString `db:"twitter"`
		Discord       sql.NullString `db:"discord"`
		Avatar        sql.NullString `db:"avatar"`
		Banner        sql.NullString `db:"banner"`
		Timestamp     int64          `db:"timestamp"`
		TwitterCreate sql.NullInt64  `db:"twitter_create"`
		EmailCreate   sql.NullInt64  `db:"email_create"`
	}
)

func newUserModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user`",
	}
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, data.Address)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, userAddressKey, userIdKey)
	return err
}

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, userIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindBySQL(ctx context.Context, sql string) (*[]User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, "timestamp")
	var resp []User
	err := m.QueryRowCtx(ctx, &resp, userIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {

		query := fmt.Sprintf("select %s from %s ", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, "timestamp")
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByAddress(ctx context.Context, address string) (*User, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, address)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userAddressKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `address` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, address); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByTwitter(ctx context.Context, address string) (*User, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserTwitterPrefix, address)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userAddressKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `twitter` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, address); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByDiscord(ctx context.Context, address string) (*User, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserDiscordPrefix, address)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userAddressKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `discord` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, address); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByUserName(ctx context.Context, username string) (*User, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserUsernamePrefix, username)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userAddressKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, username); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByNNS(ctx context.Context, nns string) (*User, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserNNSPrefix, nns)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userAddressKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `nns` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, nns); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, data.Address)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?,?)", m.table, userRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Username, data.Bio, data.Address, data.NNS, data.Email, data.Twitter, data.Discord, data.Avatar, data.Banner, data.Timestamp, data.TwitterCreate, data.EmailCreate)
	}, userAddressKey, userIdKey)
	return ret, err
}

func (m *defaultUserModel) Update(ctx context.Context, newData *User) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, data.Address)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.Username, newData.Bio, newData.Address, newData.NNS, newData.Email, newData.Twitter, newData.Discord, newData.Avatar, newData.Banner, newData.Timestamp, newData.TwitterCreate, newData.EmailCreate, newData.Id)
	}, userAddressKey, userIdKey)
	return err
}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserModel) tableName() string {
	return m.table
}
