package main

import (
	"flag"
	"fmt"
	"magaOasis/internal/config"
	"magaOasis/internal/handler"
	"magaOasis/internal/svc"
	//"magaOasis/ws/email"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/megaoasis-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
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
	//group.Add(internal.Server{Ctx: ctx,Config: c})
	//group.Start()
}


