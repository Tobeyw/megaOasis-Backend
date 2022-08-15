package home

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"magaOasis/internal/config"
	"math"
	"math/big"
)

func (me *T) GetAuctionSuccess(currentTime int64, intervalTime int64, cfg config.Config) {

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
		Index:      "getBidSuccess",
		Sort:       bson.M{},
		Filter: bson.M{
			"auctionType": 2,
			"bidAmount":   bson.M{"$gt": 0},
			"$and": []interface{}{
				bson.M{"deadline": bson.M{"$lt": currentTime}},
				bson.M{"deadline": bson.M{"$gte": beforeTime}},
			}},
		Query: []string{},
	}, ret)
	if err != nil {
		fmt.Println("get auction success event failed:", err)
	}
	if count > 0 {
		for _, item := range r1 {
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)

			auctor := item["auctor"].(string)
			bidder := item["bidder"].(string)
			auctionAsset := item["auctionAsset"].(string)
			//bidAmount := item["bidAmount"].
			item["bidAmount"].(primitive.Decimal128).String()
			bidAmount, _, err := item["bidAmount"].(primitive.Decimal128).BigInt()
			if err != nil {
				fmt.Println("bigInt convert err:", err)
			}
			ibf := new(big.Float).SetInt(bidAmount)

			symbol, decimal, err := me.GetAssetSymbol(auctionAsset)
			if err != nil {
				fmt.Println("GetAuctionAssetInfo err : ", err)
			}
			convertAmount := new(big.Float).Quo(ibf, big.NewFloat(math.Pow(10, float64(decimal))))
			nftName, err := me.GetNFTName(asset, tokenid)
			if err != nil {
				fmt.Println("GetNftName err:", err)
			}
			// 通知卖家拍卖成功
			auctorEmail, err := GetEmail(auctor)
			if err != nil {
				fmt.Println("GetAuctorEmail err:", err)
			}
			if auctorEmail != "" {
				auctorSubject := "Congratulations, your item auctioned!"
				auctorBody := "You successfully auctioned " + nftName + " for " + convertAmount.String() + symbol + " on MegaOasis."
				SendEmailOutLook(cfg, auctorSubject, auctorBody, auctorEmail)
				fmt.Println("sende email to " + auctorEmail + "for auction success")
			}

			//通知买家拍卖成功
			bidderEmail, err := GetEmail(bidder)
			if err != nil {
				fmt.Println("GetBidderEmail err:", err)
			}
			if bidderEmail != "" {
				bidderSubject := "Congratulations, you bought an item at auction!"
				bidderBody := "You successfully bought " + nftName + " for " + convertAmount.String() + symbol + " at auction on MegaOasis."
				SendEmailOutLook(cfg, bidderSubject, bidderBody, bidderEmail)
				fmt.Println("sende email to " + bidderEmail + "for bid success")
			}

		}
	}

}
