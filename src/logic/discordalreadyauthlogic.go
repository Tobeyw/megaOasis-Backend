package logic

import (
	"context"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiscordAlreadyAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscordAlreadyAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscordAlreadyAuthLogic {
	return &DiscordAlreadyAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscordAlreadyAuthLogic) DiscordAlreadyAuth() (resp *types.Response, err error) {
	return &types.Response{
		Code:    32001,
		Message: "That discord account is already taken. Please try another one.",
	}, nil
}
