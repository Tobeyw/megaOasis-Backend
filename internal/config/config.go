package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	User zrpc.RpcClientConf
	DataSource string
	Table string
	Cache cache.CacheConf
	Email struct{
		Account string
		Passwd string
		Host string
	}
	MongoDB string
	DB string
}
