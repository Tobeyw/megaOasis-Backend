package logic

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"magaOasis/common/consts"
	"magaOasis/lib/type/nullstring"
	"net/http"
	"net/url"
	"os"
	"strings"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindTwitterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindTwitterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindTwitterLogic {
	return &BindTwitterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindTwitterLogic) BindTwitter(req *types.TwitterAccessToken) (resp *types.Response, err error) {

	verifyResult, err := isVerify(req.Signature)
	if err != nil {
		return &types.Response{Message: "signature verify failed"}, err
	}
	if verifyResult {
		accessToken, err := GetAccessTokenFromCode(req.Code)
		if err != nil {
			return &types.Response{Code: 32001, Message: "GetAccessTokenFromCode failed"}, err
		}
		userName, err := GetUserInfoTwitter(accessToken)
		if err != nil {
			return &types.Response{Code: 32001, Message: "GetUserInfoTwitter failed"}, err
		}

		user, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, req.Address)
		if err != nil {
			return &types.Response{Code: 32001, Message: "FindUserByAddress failed"}, err
		}
		user.Twitter = sql.NullString{userName, nullstring.IsNull(userName)}

		err = l.svcCtx.UserModel.Update(l.ctx, user)
		if err != nil {
			return &types.Response{Code: 32001, Message: "UserInfoUpdate failed"}, err
		}
	} else {
		return &types.Response{Message: "signature verify false"}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "success",
	}, nil
}

func GetAccessTokenFromCode(code string) (accessToken string, err error) {
	postData := url.Values{}
	postData.Add("code", code)
	postData.Add("client_id", consts.TwitterClientID)
	postData.Add("code_verifier", consts.TwitterCodeChallenge)
	postData.Add("grant_type", consts.TwitterAccessTokenGrantType)
	rt := os.ExpandEnv("${RUNTIME}")
	if rt == "test" {
		postData.Add("redirect_uri", consts.TwitterRedirectURITest)
	} else {
		postData.Add("redirect_uri", consts.TwitterRedirectURIMain)
	}

	client := &http.Client{}

	url := consts.TwitterAccessTokenEndpoint
	reqest, err := http.NewRequest("POST", url, strings.NewReader(postData.Encode()))

	authstr := consts.TwitterClientID + ":" + consts.TwitterClientScret

	basic := base64.StdEncoding.EncodeToString([]byte(authstr))
	fmt.Println("basic:", basic)

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//reqest.Header.Set("Authorization","Basic "+basic)
	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		//log.Errorf("reader error:%v", err)
		return "", err
	}
	defer response.Body.Close()
	reader := response.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	if err1 := json.Unmarshal(body, &data); err1 != nil {
		return "", err
	}
	token := data["access_token"]
	if token != nil {
		accessToken = data["access_token"].(string)
	}
	fmt.Println("accessToken:", accessToken)
	return accessToken, nil
}

func GetUserInfoTwitter(accessToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, consts.TwitterGetUserInfoEndpoint, nil)
	if err != nil {
		//log.Errorf("make request error:%v", err)
		return "", err
	}
	var bearer = "Bearer " + accessToken
	req.Header.Add("Authorization", bearer)

	fmt.Println("GetUserInfoTwitter para:", req)
	resp, err := client.Do(req)
	if err != nil {
		//log.Errorf("send request error:%v", err)
		return "", err
	}
	defer resp.Body.Close()
	reader := resp.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("GetUserInfoTwitter err:", err)
		return "", err
	}

	fmt.Println("body: ", string(body))
	var data map[string]interface{}
	if err1 := json.Unmarshal(body, &data); err1 != nil {
		return "", err
	}
	username := ""
	if data["data"] != nil {
		uname := data["data"].(map[string]interface{})["username"]

		if uname != nil {
			username = data["data"].(map[string]interface{})["username"].(string)
		}
	}

	return username, nil
}
