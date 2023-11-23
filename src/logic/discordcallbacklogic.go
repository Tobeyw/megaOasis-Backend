package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"magaOasis/common/consts"
	"magaOasis/lib/type/nullstring"
	"magaOasis/model/user"
	"magaOasis/src/svc"
	"magaOasis/src/types"
	"net/http"
	"os"
	"time"
)

type DiscordCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscordCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscordCallbackLogic {
	return &DiscordCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscordCallbackLogic) DiscordCallback(req *types.CallbackDiscordParam, w http.ResponseWriter, r *http.Request) {
	rt := os.ExpandEnv("${RUNTIME}")
	url := consts.FrontEndRedirectUrlMain

	if rt == "test" {
		url = consts.FrontEndRedirectUrlTest
	}

	code := req.Code
	address := req.State

	accessToken, err := GetDiscordAccessTokenFromCode(code)
	if err != nil || accessToken == "" {
		log.Println("GetAccessTokenFromCode failed ", err)
		//return &types.Response{"GetAccessTokenFromCode failed"}, err
	}
	userName, err := GetUserInfoDiscord(accessToken)
	if err != nil || userName == "" {
		log.Println("GetUserInfoTwitter failed ", err)
		//return &types.Response{"GetUserInfoTwitter failed"}, err
	}

	//查看该twitter是否验证过
	getTwitter, err := l.svcCtx.UserModel.FindOneByDiscord(l.ctx, userName)
	fmt.Println("twitter:", getTwitter)
	if getTwitter == nil {
		getuser, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, address)
		if err != nil && err.Error() != "sql: no rows in result set" {
			log.Println("GetUserInfo failed ", err)
		}

		fmt.Println(getuser, err)

		if getuser == nil {
			// add
			_, err := l.svcCtx.UserModel.Insert(l.ctx, &user.User{
				Username:      sql.NullString{"", nullstring.IsNull("")},
				Address:       address,
				Bio:           sql.NullString{"", nullstring.IsNull("")},
				Email:         sql.NullString{"", nullstring.IsNull("")},
				Twitter:       sql.NullString{"", nullstring.IsNull("")},
				Discord:       sql.NullString{userName, nullstring.IsNull(userName)},
				Avatar:        sql.NullString{"", nullstring.IsNull("")},
				Banner:        sql.NullString{"", nullstring.IsNull("")},
				Timestamp:     time.Now().UnixMilli(),
				TwitterCreate: sql.NullInt64{time.Now().UnixMilli(), true},
				EmailCreate:   sql.NullInt64{int64(0), nullstring.IsZero(int64(0))},
			})

			if err != nil {
				log.Println("UserInfoAdd failed ", err)
			}

		} else {
			//update
			getuser.Discord = sql.NullString{userName, nullstring.IsNull(userName)}
			//getuser.TwitterCreate = sql.NullInt64{time.Now().UnixMilli(), true}
			err = l.svcCtx.UserModel.Update(l.ctx, getuser)
			if err != nil {
				log.Println("UserInfoUpdate failed ", err)

			}
		}

		fmt.Println(userName)
		//log.Fatal(getuser)

		http.Redirect(w, r, url, http.StatusFound)
	} else {
		errURL := consts.DiscordErrorPage
		http.Redirect(w, r, errURL, http.StatusFound)
	}

}
