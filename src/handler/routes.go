// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"magaOasis/src/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/profile/upload",
				Handler: UploadUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/profile/get",
				Handler: GetUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/profile/twitter",
				Handler: AuthTwitterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/twitter/callback",
				Handler: TwitterCallbackHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/unbindTwitter",
				Handler: UnbindTwitterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/twitter/bindTwitter",
				Handler: BindTwitterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/twitter/error",
				Handler: TwitterAlreadyAuthHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/airdrop",
				Handler: AirdropInfoHandler(serverCtx),
			},
		},
	)
}
