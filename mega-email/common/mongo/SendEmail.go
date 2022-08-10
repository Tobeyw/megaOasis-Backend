package neo

import (
	"fmt"
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

func SendEmail(subject, body string, to string) {
	user := "tobey1024@126.com"
	pwd := "HPIPKBZYVJWTVZXM"
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
