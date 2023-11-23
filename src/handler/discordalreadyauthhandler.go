package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"magaOasis/src/logic"
	"magaOasis/src/svc"
)

func DiscordAlreadyAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewDiscordAlreadyAuthLogic(r.Context(), svcCtx)
		resp, err := l.DiscordAlreadyAuth()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
