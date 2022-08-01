package email

import (
	"context"
	"fmt"
	neo "magaOasis/common/mongo"
	"magaOasis/common/nftEvent"
	"magaOasis/internal/logic"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"net/smtp"
	"strings"
)

func sendToMail(user,password,host,to,subject,body,mailtype string) error {
	hp := strings.Split(host,":")
	auth := smtp.PlainAuth("",user,password,hp[0])
	var content_type string
	if mailtype =="html" {
		content_type = "Content_Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content_Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To:" + to +"\r\nFrom: " + user + "<"+
		user + ">\r\nSubject: "+ subject + "\r\n" +
		content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to,";")
	err := smtp.SendMail(host,auth,user,send_to,msg)
	return err
}

func SendEmail(subject,body string,to string)  {
	user := "tobey1024@126.com"
	pwd := "HPIPKBZYVJWTVZXM"
	host := "smtp.126.com:25"
	//to := "1832541104@qq.com"//可以用;隔开发送多个
	fmt.Println("send email")
	err := sendToMail(user,pwd,host,to,subject,body,"html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

//获取megaoasis 事件
func  getEvent(fura *neo.T,svcCtx *svc.ServiceContext)  {
	fmt.Println("event")

	result, err :=fura.GetAddressCount()
	if err != nil{
		fmt.Println("Error",err)
	}
	event := result["nftEvent"].(nftEvent.T).Val()
	if event == nftEvent.Sold_Success.Val(){
		name:=result["name"].(string)
		amount:=result["convertAmount"].(string)
		symbol:=result["symbol"].(string)
		address:=types.Address{
			Address: result["owner"].(string),
		}
		to,err :=GetEmail(address,svcCtx)
		if err!=nil{
			fmt.Println("Error:",err)
			return
		}
		fmt.Println("email: ",to)
		title:="Congratulations, your item sold!"
		body:="You successfully sold "+ name +" for "+ amount +" "+ symbol+" on MegaOasis."
		SendEmail(title,body,to)

	}else if event == nftEvent.Receive_Offer.Val() {
		name:=result["name"].(string)
		amount:=result["convertAmount"].(string)
		symbol:=result["symbol"].(string)
		address:=types.Address{
			Address:result["originOwner"].(string),
		}
		to,err :=GetEmail(address,svcCtx)
		if err!=nil{
			fmt.Println("Error:",err)
			return
		}
		title:="Someone made an offer on your item!"
		body:="You have an offer of "+ amount +" "+ symbol +" for "+ name +" on MegaOasis."
		SendEmail(title,body,to)

	}else if event == nftEvent.Accept_Offer.Val() {
		name:=result["name"].(string)
		amount:=result["convertAmount"].(string)
		symbol:=result["symbol"].(string)

		address:=types.Address{
			Address:result["offerer"].(string),
		}
		to,err :=GetEmail(address,svcCtx)
		if err!=nil{
			fmt.Println("Error:",err)
			return
		}
		title:="Congratulations, your offer was accepted!"
		body:="Your offer of "+amount+" "+symbol+" for "+ name +" was accepted on MegaOasis."
		SendEmail(title,body,to)
	}
}

func GetEmail(req types.Address,svcCtx *svc.ServiceContext)  (string,error) {
	l := logic.NewGetUserLogic(context.Background(), svcCtx)
	resp, err := l.GetUser(&req)
	if err !=nil{
		return "",nil
	}
	email :=resp.Email
	return  email,nil
}

func  getCronEvent(fura *neo.T,svcCtx *svc.ServiceContext)  {
	fmt.Println("cron event")
	result,count, err :=fura.GetExpiredNFT()
	if err != nil{
		fmt.Println("Error",err)
	}
	if count >0 {
		for _, item := range result {
			address:=types.Address{
				Address:item["auctor"].(string),
			}
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)

			name ,err :=fura.GetNFTName(asset,tokenid)
			if err != nil{
				fmt.Println("Error:",err)
				return
			}
			to,err :=GetEmail(address,svcCtx)
			if err!=nil{
				fmt.Println("Error:",err)
				return
			}
			title := "Attention! Your item list expired!"
			body := ":Your listing of "+name+" has expired on MegaOasis."
			SendEmail(title,body,to)

		}
	}
	
}



