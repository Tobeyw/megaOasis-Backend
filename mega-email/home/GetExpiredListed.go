package home

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"magaOasis/internal/config"
)

func (me *T) GetExpiredListed(currentTime int64, intervalTime int64, cfg config.Config) {

	beforeTime := currentTime - intervalTime

	fmt.Println("time:", beforeTime, currentTime)
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
		Collection: "Market",
		Index:      "getExpiredListed",
		Sort:       bson.M{},
		Filter: bson.M{
			"bidAmount": 0,
			"$and": []interface{}{
				bson.M{"deadline": bson.M{"$lt": currentTime}},
				bson.M{"deadline": bson.M{"$gte": beforeTime}},
			}},
		Query: []string{},
	}, ret)
	if err != nil {
		fmt.Println("get expired list NFT failed:", err)
	}
	if count > 0 {
		for _, item := range r1 {
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)
			auctor := item["auctor"].(string)

			nftName, err := me.GetNFTName(asset, tokenid)
			if err != nil {
				fmt.Println("GetNftName err:", err)
			}
			// 通知卖家拍卖成功
			auctorEmail, err := GetEmail(auctor)
			if err != nil {
				fmt.Println("Get lister Email err:", err)
			}
			if auctorEmail != "" {
				auctorSubject := "Attention! Your item list expired!"
				auctorBody := "Your listing of " + nftName + " has expired on MegaOasis."
				SendEmailOutLook(cfg, auctorSubject, auctorBody, auctorEmail)
				fmt.Println("sende email to " + auctorEmail + "for NFT listed expired")
			}

		}
	}

}
