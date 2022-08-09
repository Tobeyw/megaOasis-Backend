package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"magaOasis/common/consts"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"net/http"
)

type TwitterCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTwitterCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TwitterCallbackLogic {
	return &TwitterCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TwitterCallbackLogic) TwitterCallback(req *types.CallbackTwitterParam, w http.ResponseWriter, r *http.Request) {

	url := consts.FrontEndRedirectUrlTest

	code := req.Code
	address := req.State

	accessToken, err := GetAccessTokenFromCode(code)
	if err != nil {
		log.Println("GetAccessTokenFromCode failed ", err)
		//return &types.Response{"GetAccessTokenFromCode failed"}, err
	}
	userName, err := GetUserInfoTwitter(accessToken)
	if err != nil {
		log.Println("GetUserInfoTwitter failed ", err)
		//return &types.Response{"GetUserInfoTwitter failed"}, err
	}

	getuser, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, address)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Println("GetUserInfo failed ", err)
	}

	fmt.Println(getuser, err)

	//if getuser == nil {
	//	// add
	//	_, err := l.svcCtx.UserModel.Insert(l.ctx, &user.User{
	//		Username:  sql.NullString{"", nullstring.IsNull("")},
	//		Address:   address,
	//		Bio:       sql.NullString{"", nullstring.IsNull("")},
	//		Email:     sql.NullString{"", nullstring.IsNull("")},
	//		Twitter:   sql.NullString{userName, nullstring.IsNull(userName)},
	//		Avatar:    sql.NullString{"", nullstring.IsNull("")},
	//		Banner:    sql.NullString{"", nullstring.IsNull("")},
	//		Timestamp: time.Now().UnixMilli(),
	//	})
	//
	//	if err != nil {
	//		log.Println("UserInfoAdd failed ", err)
	//	}
	//
	//} else {
	//	//update
	//	getuser.Twitter = sql.NullString{userName, nullstring.IsNull(userName)}
	//	err = l.svcCtx.UserModel.Update(l.ctx, getuser)
	//	if err != nil {
	//		log.Println("UserInfoUpdate failed ", err)
	//		//return &types.Response{"UserInfoUpdate failed"}, err
	//
	//	}
	//}

	fmt.Println(userName)
	//log.Fatal(getuser)

	http.Redirect(w, r, url, http.StatusFound)

}
