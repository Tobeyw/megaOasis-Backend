package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"magaOasis/src/config"
	"magaOasis/src/handler"
	"magaOasis/src/svc"
)

var configFile = flag.String("f", "etc/megaoasis-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.RestConf.MaxBytes = -1
	c.RestConf.Timeout = 60000 //设置传输时间 1min
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	//cd,dbonline :=intializeMongoOnlineClient(c,context.TODO())
	//fura := &neo.T{
	//	Db_online: dbonline,
	//	C_online: cd,
	//}
	//fura.GetAddressCount()
	//c1 := cron.New()
	//err := c1.AddFunc("@every 5s", func() {
	//
	//	go fura.GetAddressCount()
	//
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}

	//c1.Start()

	ctx := svc.NewServiceContext(c)

	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

	// add wss server
	//group := service.NewServiceGroup()
	//defer group.Stop()
	//group.Add(server)
	//group.Add(src.Server{Ctx: ctx,Config: c})
	//group.Start()
}
