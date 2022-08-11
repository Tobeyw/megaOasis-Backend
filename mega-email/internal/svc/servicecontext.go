package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"magaOasis/internal/config"
	"magaOasis/model"

	"os"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel model.MessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	rt := os.ExpandEnv("${RUNTIME}")

	if rt == "test" {
		return &ServiceContext{
			Config:       c,
			MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSourceTest), c.CacheTest),
		}
	} else if rt == "main" {
		return &ServiceContext{
			Config:       c,
			MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSourceMain), c.CacheMain),
		}
	} else if rt == "dev" {
		return &ServiceContext{
			Config:       c,
			MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSourceDev), c.CacheDev),
		}
	} else {
		return &ServiceContext{
			Config:       c,
			MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSourceMain), c.CacheMain),
		}
	}
}
