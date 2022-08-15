package home

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	gomail "gopkg.in/mail.v2"
	"io/ioutil"
	"log"
	"magaOasis/internal/config"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

func GetEmail(address string) (string, error) {
	rt := os.ExpandEnv("${RUNTIME}")
	url := "https://megaoasis.ngd.network:8893/profile/get?address=" + address
	if rt == "test" {
		url = "https://megaoasis.ngd.network:8889/profile/get?address=" + address
	} else if rt == "dev" {
		url = "http://localhost:8889/profile/get?address=" + address
	}
	fmt.Println("getUser :" + url)

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

	fmt.Println(string(body))
	return email, nil

}

func SendEmailOutLook(cfg config.Config, subject, body string, to string) {
	fmt.Println(subject, to)
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
	} else {
		fmt.Println("send email to " + to + "failed")
	}

}

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

func (me *T) GetNFTName(asset string, tokenid string) (string, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	r1, count, err := me.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Nep11Properties",
		Index:      "Nep11Properties",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": asset, "tokenid": tokenid},
		Query:      []string{},
		Limit:      10,
		Skip:       0,
	}, ret)
	if err != nil {
		return "", err
	}
	nftname := ""
	if count == int64(1) {
		properties := r1[0]["properties"].(string)
		if properties != "" {
			var data map[string]interface{}
			if err1 := json.Unmarshal([]byte(properties), &data); err1 == nil {
				name, ok := data["name"].(string)
				if ok {
					nftname = name
				} else {
					nftname = ""
				}
			} else {
				return "", err
			}
		}
	}
	return nftname, nil
}

func (me *T) GetAssetSymbol(asset string) (string, int32, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	r1, count, err := me.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Asset",
		Index:      "Asset",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": asset},
		Query:      []string{},
		Limit:      10,
		Skip:       0,
	}, ret)
	if err != nil {
		return "", 0, err
	}
	symbol := ""
	decimals := int32(0)
	if count == int64(1) {
		decimals = r1[0]["decimals"].(int32)
		symbol = r1[0]["symbol"].(string)
	}

	return symbol, decimals, nil
}

func (me *T) GetOwner(nonce int64) (string, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	r1, count, err := me.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "MarketNotification",
		Index:      "MarketNotification",
		Sort:       bson.M{},
		Filter:     bson.M{"eventname": "Auction", "nonce": nonce},
		Query:      []string{},
		Limit:      10,
		Skip:       0,
	}, ret)
	if err != nil {
		return "", err
	}
	var user string
	if count == int64(1) {
		user = r1[0]["user"].(string)
	}
	return user, nil
}
