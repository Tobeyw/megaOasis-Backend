package logic

import (
	"context"
	"magaOasis/common/consts"

	"magaOasis/internal/svc"
	"magaOasis/internal/types"

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
	entryUrl := consts.TwitterAuthorizeEndpoint + "?" + "response_type=" + consts.TwitterResponseType + "&client_id=" + consts.TwitterClientID + "&scope=" + consts.TwitterScope + "&state=" + req.Address + "&redirect_uri=" + consts.TwitterRedirectURI + "&code_challenge=" + consts.TwitterCodeChallenge + "&code_challenge_method=" + consts.TwitterCodeChallengeMethod
	return &types.LoginTwitterResponse{
		Url: entryUrl,
	}, nil

	return
}
