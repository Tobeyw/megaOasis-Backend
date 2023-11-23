package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"magaOasis/src/logic"
	"magaOasis/src/svc"
	"magaOasis/src/types"
)

func AuthDiscordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var req types.Address
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewAuthDiscordLogic(r.Context(), svcCtx)
		resp, err := l.AuthDiscord(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
