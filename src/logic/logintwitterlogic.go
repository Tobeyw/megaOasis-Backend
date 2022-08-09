package logic

import (
	"context"
	"magaOasis/common/consts"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginTwitterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginTwitterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginTwitterLogic {
	return &LoginTwitterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginTwitterLogic) LoginTwitter() (resp *types.LoginTwitterResponse, err error) {
	entryUrl := consts.TwitterAuthorizeEndpoint + "?" + "response_type=" + consts.TwitterResponseType + "&client_id=" + consts.TwitterClientID + "&scope=" + consts.TwitterScope + "&state=" + "state" + "&redirect_uri=" + consts.TwitterRedirectURI + "&code_challenge=" + consts.TwitterCodeChallenge + "&code_challenge_method=" + consts.TwitterCodeChallengeMethod

	return &types.LoginTwitterResponse{
		Url: entryUrl,
	}, nil
}
