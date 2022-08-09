package logic

import (
	"context"
	"net/http"

	"magaOasis/internal/svc"
	"magaOasis/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

	//url := "https://www.baidu.com/"
	//
	//code := req.Code
	//address := req.State
	//
	//accessToken, err := GetAccessTokenFromCode(code)
	//if err != nil {
	//	//return &types.Response{"GetAccessTokenFromCode failed"}, err
	//}
	//userName, err := GetUserInfoTwitter(accessToken)
	//if err != nil {
	//	//return &types.Response{"GetUserInfoTwitter failed"}, err
	//}
	//
	//user, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, address)
	//if err != nil {
	//	//return &types.Response{"FindUserByAddress failed"}, err
	//}
	//user.Twitter = sql.NullString{userName, nullstring.IsNull(userName)}
	//
	//err = l.svcCtx.UserModel.Update(l.ctx, user)
	//if err != nil {
	//	//return &types.Response{"UserInfoUpdate failed"}, err
	//}
	//
	//log.Fatal(userName)
	//
	//http.Redirect(w, r, url, http.StatusFound)

}
