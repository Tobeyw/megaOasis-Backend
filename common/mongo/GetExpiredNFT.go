package neo

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (me *T) GetExpiredNFT() ([]map[string]interface{},int64,error) {
	currentTime := time.Now().UnixNano() / 1e6
    before5min := currentTime - 5*60*1000
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
		Index:      "Market",
		Sort:       bson.M{},
		Filter:     bson.M{"amount":1,"bidAmount":0,"$and":[]interface{}{bson.M{"deadline":bson.M{"$lte":currentTime}},bson.M{"deadline":bson.M{"$gt":before5min}}}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return nil,0,err
	}

	return r1, count,nil
}




