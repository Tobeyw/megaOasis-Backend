package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"magaOasis/internal/logic"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
)

func AuthTwitterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Address
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewAuthTwitterLogic(r.Context(), svcCtx)
		resp, err := l.AuthTwitter(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
