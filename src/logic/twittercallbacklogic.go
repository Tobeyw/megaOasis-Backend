package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"magaOasis/src/svc"
	"magaOasis/src/types"
)

type TwitterCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTwitterCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TwitterCallbackLogic {
	return &TwitterCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TwitterCallbackLogic) TwitterCallback(req *types.CallbackTwitterParam) (resp *types.Response, err error) {
	code := req.Code

	return &types.Response{Message: code}, nil
}
