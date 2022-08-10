package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"magaOasis/model/user"
	"magaOasis/src/config"
	"os"
)

type ServiceContext struct {
	Config    config.Config
	UserModel user.UserModel
	MysqlConn sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	rt := os.ExpandEnv("${RUNTIME}")
	fmt.Println(rt)
	if rt == "test" {
		return &ServiceContext{
			Config:    c,
			UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSourceTest), c.Cache),
		}
	} else if rt == "main" {
		return &ServiceContext{
			Config:    c,
			UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSourceMain), c.Cache),
		}
	} else {
		return &ServiceContext{
			Config:    c,
			UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSourceMain), c.Cache),
		}
	}

}
