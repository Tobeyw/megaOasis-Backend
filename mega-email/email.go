package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gomail "gopkg.in/mail.v2"
	"io/ioutil"
	"log"
	"magaOasis/common/nftEvent"
	"magaOasis/home"
	"os"

	"magaOasis/internal/config"
	"magaOasis/internal/handler"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"math"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/email-api.yaml", "the config file")

func main() {
	log.Println("YOUR ENV IS %s", os.ExpandEnv("${RUNTIME}"))
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	//==============
	cd, dbonline := intializeMongoOnlineClient(c, context.TODO())
	me := home.T{
		Db_online: dbonline,
		C_online:  cd,
	}

	conn, err := me.GetCollection(struct{ Collection string }{Collection: "MarketNotification"})
	if err != nil {
		fmt.Println("conn :", err)
	}
	cs, err := conn.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		//return nil,err
		fmt.Println("watch:", err)
	}

	fmt.Println("watching....")
	for cs.Next(context.TODO()) {
		var changeEvent map[string]interface{}
		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}
		eventItem := changeEvent["fullDocument"].(map[string]interface{})
		eventname := eventItem["eventname"]
		asset := eventItem["asset"].(string)
		tokenid := eventItem["tokenid"].(string)
		fmt.Println(eventname)
		if eventname == "Claim" {
			nonce := eventItem["nonce"].(int64)
			extendData := eventItem["extendData"].(string)
			var dat map[string]interface{}
			if err31 := json.Unmarshal([]byte(extendData), &dat); err31 == nil {
				eventItem["auctionAsset"] = dat["auctionAsset"].(string)
				eventItem["bidAmount"] = dat["bidAmount"].(string)
				eventItem["auctionType"] = dat["auctionType"].(string)
				if dat["auctionType"].(string) == "1" {
					eventItem["nftEvent"] = nftEvent.Sold_Success
				}
				symbol, decimals, err := me.GetAssetSymbol(dat["auctionAsset"].(string))
				if err != nil {
					//return nil, err
				}
				eventItem["symbol"] = symbol
				eventItem["decimals"] = decimals

				offerAmount, err := strconv.ParseInt(dat["bidAmount"].(string), 10, 64)
				if err != nil {
					//return nil, err
				}

				convertAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(offerAmount)/math.Pow(10, float64(decimals))), 64)
				if err != nil {
					//return nil, err
				}
				eventItem["symbol"] = symbol
				eventItem["convertAmount"] = strconv.FormatFloat(convertAmount, 'f', 2, 64)

				owner, err := me.GetOwner(nonce)
				if err != nil {
					//return nil, err
				}
				eventItem["owner"] = owner
			}

		} else if eventname == "Offer" {
			extendData := eventItem["extendData"].(string)
			var dat map[string]interface{}
			if err31 := json.Unmarshal([]byte(extendData), &dat); err31 == nil {
				eventItem["originOwner"] = dat["originOwner"].(string)
				eventItem["offerAsset"] = dat["offerAsset"].(string)
				eventItem["offerAmount"] = dat["offerAmount"].(string)
				eventItem["deadline"] = dat["deadline"].(string)
				eventItem["nftEvent"] = nftEvent.Receive_Offer
				symbol, decimals, err := me.GetAssetSymbol(dat["offerAsset"].(string))
				if err != nil {
					//return nil, err
				}
				offerAmount, err := strconv.ParseInt(dat["offerAmount"].(string), 10, 64)
				if err != nil {
					//return nil, err
				}
				convertAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(offerAmount)/math.Pow(10, float64(decimals))), 64)
				if err != nil {
					//return nil, err
				}
				eventItem["symbol"] = symbol
				eventItem["convertAmount"] = strconv.FormatFloat(convertAmount, 'f', 2, 64)
			}

		} else if eventname == "CompleteOffer " {
			extendData := eventItem["extendData"].(string)
			var dat map[string]interface{}
			if err31 := json.Unmarshal([]byte(extendData), &dat); err31 == nil {
				eventItem["offerer"] = dat["offerer"].(string)
				eventItem["offerAsset"] = dat["offerAsset"].(string)
				eventItem["offerAmount"] = dat["offerAmount"].(string)
				eventItem["deadline"] = dat["deadline"].(string)
				eventItem["nftEvent"] = nftEvent.Accept_Offer
				symbol, decimals, err := me.GetAssetSymbol(dat["offerAsset"].(string))
				if err != nil {
					//return nil, err
				}
				eventItem["symbol"] = symbol
				eventItem["decimals"] = decimals

				offerAmount, err := strconv.ParseInt(dat["offerAmount"].(string), 10, 64)
				if err != nil {
					//return nil, err
				}
				convertAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(offerAmount)/math.Pow(10, float64(decimals))), 64)
				if err != nil {
					//return nil, err
				}
				eventItem["symbol"] = symbol
				eventItem["convertAmount"] = strconv.FormatFloat(convertAmount, 'f', 2, 64)
			}

		} else {
			//return nil, nil
		}

		nftname, err := me.GetNFTName(asset, tokenid)
		if err != nil {
			//return nil, err
		}
		eventItem["name"] = nftname
		SendEmailByEvent(c, ctx, eventItem)

	}

	//==============
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}

func intializeMongoOnlineClient(cfg config.Config, ctx context.Context) (*mongo.Client, string) {
	rt := os.ExpandEnv("${RUNTIME}")
	//默认main
	clientOptions := options.Client().ApplyURI(cfg.MongoDBMain)
	dbOnline := cfg.DBMain
	if rt == "test" {
		clientOptions = options.Client().ApplyURI(cfg.MongoDBTest)
		dbOnline = cfg.DBTest
	}

	clientOptions.SetMaxPoolSize(20)
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	return co, dbOnline
}

func SendEmailByEvent(cfg config.Config, svcCtx *svc.ServiceContext, result map[string]interface{}) {
	if result["nftEvent"] != nil {
		event := result["nftEvent"].(nftEvent.T).Val()
		if event == nftEvent.Sold_Success.Val() {
			name := result["name"].(string)
			amount := result["convertAmount"].(string)
			symbol := result["symbol"].(string)
			address := types.Address{
				Address: result["owner"].(string),
			}
			to, err := GetEmail(address, svcCtx)
			if err != nil {
				fmt.Println("Error:", err)
				//return
			}
			fmt.Println("email: ", to)
			title := "Congratulations, your item sold!"
			body := "You successfully sold " + name + " for " + amount + " " + symbol + " on MegaOasis."
			SendEmailOutLook(cfg, title, body, to)
			//return title,body,to

		} else if event == nftEvent.Receive_Offer.Val() {
			name := result["name"].(string)
			amount := result["convertAmount"].(string)
			symbol := result["symbol"].(string)
			address := types.Address{
				Address: result["originOwner"].(string),
			}
			to, err := GetEmail(address, svcCtx)
			if err != nil {
				fmt.Println("Error:", err)
				//return
			}
			title := "Someone made an offer on your item!"
			body := "You have an offer of " + amount + " " + symbol + " for " + name + " on MegaOasis."
			SendEmailOutLook(cfg, title, body, to)
			//return title,body,to

		} else if event == nftEvent.Accept_Offer.Val() {
			name := result["name"].(string)
			amount := result["convertAmount"].(string)
			symbol := result["symbol"].(string)

			address := types.Address{
				Address: result["offerer"].(string),
			}
			to, err := GetEmail(address, svcCtx)
			if err != nil {
				fmt.Println("Error:", err)
				//return
			}
			title := "Congratulations, your offer was accepted!"
			body := "Your offer of " + amount + " " + symbol + " for " + name + " was accepted on MegaOasis."
			SendEmailOutLook(cfg, title, body, to)
			//return title,body,to

		}
	}

}

func SendEmailOutLook(cfg config.Config, subject, body string, to string) {
	if to != "" {
		m := gomail.NewMessage()               // 声明一封邮件对象
		m.SetHeader("From", cfg.Email.Account) // 发件人
		m.SetHeader("To", to)                  // 收件人
		m.SetHeader("Subject", subject)        // 邮件主题
		m.SetBody("text/plain", body)          // 邮件内容

		// host 是提供邮件的服务器，port是服务器端口，username 是发送邮件的账号, password是发送邮件的密码
		d := gomail.NewDialer(cfg.Email.Host, cfg.Email.Port, cfg.Email.Account, cfg.Email.Passwd)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 配置tls，跳过验证
		if err := d.DialAndSend(m); err != nil {
			log.Fatalln("msg", "try send a mail failed", "err", err)
		} else {
			fmt.Println("send email to " + to)
		}
	}

}

func GetEmail(req types.Address, svcCtx *svc.ServiceContext) (string, error) {
	rt := os.ExpandEnv("${RUNTIME}")
	url := "https://megaoasis.ngd.network:8893/profile/get?address=" + req.Address
	if rt == "test" {
		url = "https://megaoasis.ngd.network:8889/profile/get?address=" + req.Address
	}

	resp, err := http.Get(url)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()
	reader := resp.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		//log.Errorf("reader error:%v", err)
		return "", err
	}

	var data map[string]interface{}
	if err1 := json.Unmarshal(body, &data); err1 != nil {
		return "", err
	}
	email := ""
	if data["email"] != nil {
		email = data["email"].(string)
	}

	return email, nil

}
