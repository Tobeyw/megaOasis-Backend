package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"magaOasis/internal/config"
	"magaOasis/model/user"
)

type ServiceContext struct {
	Config    config.Config
	UserModel user.UserModel
	MysqlConn sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSource), c.Cache),
	}
}
