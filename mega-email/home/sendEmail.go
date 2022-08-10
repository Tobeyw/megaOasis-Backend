package home

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
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
