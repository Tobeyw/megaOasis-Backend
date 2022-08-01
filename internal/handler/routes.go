// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"magaOasis/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/profile/edituser",
				Handler: EditUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/upload",
				Handler: UploadUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/testuser",
				Handler: TestUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/profile/get",
				Handler: GetUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/image",
				Handler: GetImageHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler:http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path != "/" {
							http.Error(w, "Not found", http.StatusNotFound)
							return
						}
						if r.Method != "GET" {
							http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
							return
						}

						http.ServeFile(w, r, "home.html")
					},
				),
			},

		},
	)
}


