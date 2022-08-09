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
	//code := req.Code

	//client := &http.Client{}
	url := "https://www.baidu.com/"
	//reqest, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	panic(err)
	//}
	//client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	//	return http.ErrUseLastResponse
	//}
	//response, _ := client.Do(reqest)

	http.Redirect(w, r, url, http.StatusFound)

}
