package neo

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type T struct {
	Db_online string
	C_online  *mongo.Client
}

func (me *T) GetCollection(args struct {
	Collection string
}) (*mongo.Collection, error) {
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	return collection, nil
}

func (me *T) QueryAll(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
	Query      []string
	Limit      int64
	Skip       int64
}, ret *json.RawMessage) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	op := options.Find()
	op.SetSort(args.Sort)
	op.SetLimit(args.Limit)
	op.SetSkip(args.Skip)
	co := options.CountOptions{}
	count, err := collection.CountDocuments(context.TODO(), args.Filter, &co)
	if err != nil {
		return nil, 0, fmt.Errorf("count documents errors:%s", err)
	}
	cursor, err := collection.Find(context.TODO(), args.Filter, op)
	defer cursor.Close(context.TODO())
	if err == mongo.ErrNoDocuments {
		return nil, 0, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, 0, fmt.Errorf("get cursor errors:%s", err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, 0, fmt.Errorf("find documents errors:%s", err)
	}
	for _, item := range results {
		if len(args.Query) == 0 {
			convert = append(convert, item)
		} else {
			temp := make(map[string]interface{})
			for _, v := range args.Query {
				temp[v] = item[v]
			}
			convert = append(convert, temp)
		}
	}
	r, err := json.Marshal(convert)
	if err != nil {
		return nil, 0, fmt.Errorf("json marshal errors:%s", err)
	}
	*ret = json.RawMessage(r)
	return convert, count, nil
}
