package email

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"magaOasis/common/nftEvent"
	"magaOasis/internal/config"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"math"
	"strconv"
)

// (map[string]interface{},error)
func (me *T) GetAddressCount(cfg config.Config, svcCtx *svc.ServiceContext) {

	c, err := me.GetCollection(struct{ Collection string }{Collection: "MarketNotification"})
	if err != nil {
		//return nil, err
	}
	cs, err := c.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		//return nil,err
	}
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
		SendEmailByEvent(cfg, svcCtx, eventItem)

	}

	//return nil,nil
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
