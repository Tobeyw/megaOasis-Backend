package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"magaOasis/internal/logic"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
)

func TwitterCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var req types.CallbackTwitterParam
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewTwitterCallbackLogic(r.Context(), svcCtx)
		err := l.TwitterCallback(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			//httpx.OkJson(w, resp)
		}
	}
}
