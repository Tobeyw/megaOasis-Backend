package email

import (
	"context"
	"crypto/tls"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"log"
	neo "magaOasis/common/mongo"
	"magaOasis/src/config"
	"magaOasis/src/logic"
	"magaOasis/src/svc"
	"magaOasis/src/types"
	"net/smtp"
	"strings"
)

func sendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content_Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content_Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To:" + to + "\r\nFrom: " + user + "<" +
		user + ">\r\nSubject: " + subject + "\r\n" +
		content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func SendEmailOutLook(cfg config.Config, subject, body string, to string) {
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

func SendEmail(subject, body string, to string) {
	user := "tobey1024@126.com"
	//pwd := "HPIPKBZYVJWTVZXM"
	pwd := ""
	host := "smtp.126.com:25"
	//to := "1832541104@qq.com"//可以用;隔开发送多个
	fmt.Println("send email")
	err := sendToMail(user, pwd, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

//获取megaoasis 事件
func getEvent(fura *T, cfg config.Config, svcCtx *svc.ServiceContext) {
	fmt.Println("event")
	fura.GetAddressCount(cfg, svcCtx)
	//result, err :=fura.GetAddressCount()
	//if err != nil{
	//	fmt.Println("Error",err)
	//}
	//if result["nftEvent"] !=nil {
	//
	//
	//event := result["nftEvent"].(nftEvent.T).Val()
	//if event == nftEvent.Sold_Success.Val(){
	//	name:=result["name"].(string)
	//	amount:=result["convertAmount"].(string)
	//	symbol:=result["symbol"].(string)
	//	address:=types.Address{
	//		Address: result["owner"].(string),
	//	}
	//	to,err :=GetEmail(address,svcCtx)
	//	if err!=nil{
	//		fmt.Println("Error:",err)
	//		//return
	//	}
	//	fmt.Println("email: ",to)
	//	title:="Congratulations, your item sold!"
	//	body:="You successfully sold "+ name +" for "+ amount +" "+ symbol+" on MegaOasis."
	//	SendEmail(title,body,to)
	//
	//}else if event == nftEvent.Receive_Offer.Val() {
	//	name:=result["name"].(string)
	//	amount:=result["convertAmount"].(string)
	//	symbol:=result["symbol"].(string)
	//	address:=types.Address{
	//		Address:result["originOwner"].(string),
	//	}
	//	to,err :=GetEmail(address,svcCtx)
	//	if err!=nil{
	//		fmt.Println("Error:",err)
	//		//return
	//	}
	//	title:="Someone made an offer on your item!"
	//	body:="You have an offer of "+ amount +" "+ symbol +" for "+ name +" on MegaOasis."
	//	SendEmail(title,body,to)
	//
	//}else if event == nftEvent.Accept_Offer.Val() {
	//	name:=result["name"].(string)
	//	amount:=result["convertAmount"].(string)
	//	symbol:=result["symbol"].(string)
	//
	//	address:=types.Address{
	//		Address:result["offerer"].(string),
	//	}
	//	to,err :=GetEmail(address,svcCtx)
	//	if err!=nil{
	//		fmt.Println("Error:",err)
	//		//return
	//	}
	//	title:="Congratulations, your offer was accepted!"
	//	body:="Your offer of "+amount+" "+symbol+" for "+ name +" was accepted on MegaOasis."
	//	SendEmail(title,body,to)
	//
	//}
	//}
}

func GetEmail(req types.Address, svcCtx *svc.ServiceContext) (string, error) {
	l := logic.NewGetUserLogic(context.Background(), svcCtx)
	resp, err := l.GetUser(&req)
	if err != nil {
		return "", nil
	}
	email := resp.Email
	return email, nil
}

func getCronEvent(fura *neo.T, svcCtx *svc.ServiceContext) {
	fmt.Println("cron event")
	result, count, err := fura.GetExpiredNFT()
	if err != nil {
		fmt.Println("Error", err)
	}
	if count > 0 {
		for _, item := range result {
			address := types.Address{
				Address: item["auctor"].(string),
			}
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)

			name, err := fura.GetNFTName(asset, tokenid)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			to, err := GetEmail(address, svcCtx)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			title := "Attention! Your item list expired!"
			body := ":Your listing of " + name + " has expired on MegaOasis."
			SendEmail(title, body, to)

		}
	}

}
