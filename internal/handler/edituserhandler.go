package handler

import (
	"magaOasis/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"magaOasis/internal/logic"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
)

func EditUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewEditUserLogic(r, svcCtx)
		resp, err := l.EditUser()

		response.Response(w,r,resp,err)
	}
}
