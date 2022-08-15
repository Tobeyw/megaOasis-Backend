package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"magaOasis/src/svc"
)

type AirdropInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAirdropInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AirdropInfoLogic {
	return &AirdropInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AirdropInfoLogic) AirdropInfo() error {

	return nil
}
