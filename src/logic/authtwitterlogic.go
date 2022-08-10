package logic

import (
	"context"
	"magaOasis/common/consts"
	"os"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthTwitterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthTwitterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthTwitterLogic {
	return &AuthTwitterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthTwitterLogic) AuthTwitter(req *types.Address) (resp *types.LoginTwitterResponse, err error) {
	rt := os.ExpandEnv("${RUNTIME}")
	redirect := consts.TwitterRedirectURIMain
	if rt == "test" {
		redirect = consts.TwitterRedirectURITest
	}
	entryUrl := consts.TwitterAuthorizeEndpoint + "?" + "response_type=" + consts.TwitterResponseType + "&client_id=" + consts.TwitterClientID + "&scope=" + consts.TwitterScope + "&state=" + req.Address + "&redirect_uri=" + redirect + "&code_challenge=" + consts.TwitterCodeChallenge + "&code_challenge_method=" + consts.TwitterCodeChallengeMethod
	return &types.LoginTwitterResponse{
		Url: entryUrl,
	}, nil
}
