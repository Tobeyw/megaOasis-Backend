package logic

import (
	"context"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TwitterAlreadyAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTwitterAlreadyAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TwitterAlreadyAuthLogic {
	return &TwitterAlreadyAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TwitterAlreadyAuthLogic) TwitterAlreadyAuth() (resp *types.Response, err error) {
	return &types.Response{
		Code:    32001,
		Message: "That twitter account is already taken. Please try another one.",
	}, nil
}
