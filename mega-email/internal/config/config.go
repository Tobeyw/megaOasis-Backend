package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	User           zrpc.RpcClientConf
	DataSourceDev  string
	DataSourceTest string
	DataSourceMain string
	Table          string
	CacheDev       cache.CacheConf
	CacheTest      cache.CacheConf
	CacheMain      cache.CacheConf
	Email          struct {
		Account string
		Passwd  string
		Host    string
		Port    int
	}
	MongoDBDev  string
	DBDev       string
	MongoDBTest string
	DBTest      string
	MongoDBMain string
	DBMain      string
}
