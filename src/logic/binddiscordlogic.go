package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"magaOasis/common/consts"
	"magaOasis/lib/type/nullstring"
	"net/http"
	"net/url"
	"strings"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindDiscordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindDiscordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindDiscordLogic {
	return &BindDiscordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindDiscordLogic) BindDiscord(req *types.DiscordAccessToken) (resp *types.Response, err error) {

	verifyResult, err := isVerify(req.Signature)
	if err != nil {
		return &types.Response{Message: "signature verify failed"}, err
	}
	if verifyResult {
		if true {
			accessToken, err := GetDiscordAccessTokenFromCode(req.Code)
			if err != nil {
				return &types.Response{Code: 32001, Message: "GetAccessTokenFromCode failed"}, err
			}
			userName, err := GetUserInfoDiscord(accessToken)
			if err != nil {
				return &types.Response{Code: 32001, Message: "GetUserInfoDiscord failed"}, err
			}

			user, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, req.Address)
			if err != nil {
				return &types.Response{Code: 32001, Message: "FindUserByAddress failed"}, err
			}
			user.Discord = sql.NullString{userName, nullstring.IsNull(userName)}

			err = l.svcCtx.UserModel.Update(l.ctx, user)
			if err != nil {
				return &types.Response{Code: 32001, Message: "UserInfoUpdate failed"}, err
			}
		} else {
			return &types.Response{Message: "signature verify false"}, nil
		}
	}

	return &types.Response{
		Code:    200,
		Message: "success",
	}, nil
}

func GetDiscordAccessTokenFromCode(code string) (accessToken string, err error) {
	postData := url.Values{}
	postData.Add("code", code)
	postData.Add("client_id", consts.DiscordClientID)
	postData.Add("client_secret", consts.DiscordClientSecret)
	postData.Add("grant_type", consts.DiscordAccessTokenGrantType)
	postData.Add("redirect_uri", consts.DiscordRedirectURI)

	client := &http.Client{}

	url := consts.DiscordAccessTokenEndpoint

	reqest, err := http.NewRequest("POST", url, strings.NewReader(postData.Encode()))
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//reqest.Header.Set("Authorization","Basic "+basic)

	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		//log.Errorf("reader errors:%v", err)
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

func GetUserInfoDiscord(accessToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, consts.DiscordGetUserInfoEndpoint, nil)
	if err != nil {
		//log.Errorf("make request errors:%v", err)
		return "", err
	}
	var bearer = "Bearer " + accessToken
	req.Header.Add("Authorization", bearer)

	fmt.Println("GetUserInfoDiscord para:", req)
	resp, err := client.Do(req)
	if err != nil {
		//log.Errorf("send request errors:%v", err)
		return "", err
	}
	defer resp.Body.Close()
	reader := resp.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("GetUserInfoDiscord err:", err)
		return "", err
	}
	fmt.Println("body: ", string(body))
	var data map[string]interface{}
	if err1 := json.Unmarshal(body, &data); err1 != nil {
		return "", err
	}
	username := data["username"].(string)

	return username, nil
}
