package neo

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"magaOasis/common/consts"
	"magaOasis/common/errors"
	"strings"
)

func (me *T) IsOwnerByNNS(nns string, address string) (bool, error) {

	nns = strings.TrimSpace(nns)
	var tokenid string
	if len(nns) > 0 && strings.HasSuffix(nns, ".neo") {
		tokenid = base64.URLEncoding.EncodeToString([]byte(nns))
	} else {
		return false, errors.New(30032, "nns invalid parameter")
	}

	message := make(json.RawMessage, 0)
	ret := &message
	_, count, err := me.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Address-Asset",
		Index:      "Address-Asset",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": consts.NNSContractHash, "address": address, "tokenid": tokenid},
		Query:      []string{},
	}, ret)
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}
