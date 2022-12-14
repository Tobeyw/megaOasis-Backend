// Code generated by goctl. DO NOT EDIT!

package message

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	messageFieldNames          = builder.RawFieldNames(&Message{})
	messageRows                = strings.Join(messageFieldNames, ",")
	messageRowsExpectAutoSet   = strings.Join(stringx.Remove(messageFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	messageRowsWithPlaceHolder = strings.Join(stringx.Remove(messageFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"
)

type (
	messageModel interface {
		Insert(ctx context.Context, data *Message) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Message, error)
		Update(ctx context.Context, newData *Message) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMessageModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Message struct {
		Id        int64          `db:"id"`
		Email     sql.NullString `db:"email"`
		Event     string         `db:"event"`
		Title     string         `db:"title"`
		Message   string         `db:"message"`
		Timestamp int64          `db:"timestamp"`
	}
)

func newMessageModel(conn sqlx.SqlConn) *defaultMessageModel {
	return &defaultMessageModel{
		conn:  conn,
		table: "`message`",
	}
}

func (m *defaultMessageModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMessageModel) FindOne(ctx context.Context, id int64) (*Message, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messageRows, m.table)
	var resp Message
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMessageModel) Insert(ctx context.Context, data *Message) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, messageRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Email, data.Event, data.Title, data.Message, data.Timestamp)
	return ret, err
}

func (m *defaultMessageModel) Update(ctx context.Context, data *Message) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, messageRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Email, data.Event, data.Title, data.Message, data.Timestamp, data.Id)
	return err
}

func (m *defaultMessageModel) tableName() string {
	return m.table
}
