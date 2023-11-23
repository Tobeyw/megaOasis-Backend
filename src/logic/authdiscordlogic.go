package logic

import (
	"context"
	"magaOasis/common/consts"
	"os"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthDiscordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthDiscordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthDiscordLogic {
	return &AuthDiscordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthDiscordLogic) AuthDiscord(req *types.Address) (resp *types.LoginTwitterResponse, err error) {
	rt := os.ExpandEnv("${RUNTIME}")
	redirect := consts.DiscordRedirectURIMain

	if rt == "test" {
		redirect = consts.DiscordRedirectURITest
	}
	entryUrl := consts.DiscordAuthorizeEndpoint + "?" + "response_type=" + consts.DiscordResponseType + "&client_id=" + consts.DiscordClientID + "&scope=" + consts.DiscordScope + "&state=" + req.Address + "&redirect_uri=" + redirect + "&prompt=" + consts.DiscordPrompt
	return &types.LoginTwitterResponse{
		Url: entryUrl,
	}, nil
}
