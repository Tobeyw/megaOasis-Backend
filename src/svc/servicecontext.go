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
		fmt.Println(c.DataSourceTest)
		return &ServiceContext{
			Config:    c,
			UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSourceTest), c.CacheTest),
		}
	} else if rt == "main" {
		fmt.Println(c.DataSourceMain)
		return &ServiceContext{
			Config:    c,
			UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSourceMain), c.CacheMain),
		}
	} else {
		fmt.Println(c.DataSourceDev)
		return &ServiceContext{
			Config:    c,
			UserModel: user.NewUserModel(sqlx.NewMysql(c.DataSourceDev), c.CacheDev),
		}
	}

}
