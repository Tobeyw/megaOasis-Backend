package logic

import (
	"context"

	"magaOasis/internal/svc"
	"magaOasis/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.Address) (resp *types.UserResp, err error) {
	res, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx,req.Address)

	if err != nil {
		return nil, err
	}

	return &types.UserResp{
		UserName: res.Username.String,
		Address: res.Address,
		Email: res.Email.String,
		Twitter: res.Twitter.String,
		Avatar: res.Avatar.String,
		Bio: res.Bio.String,
		Banner: res.Banner.String,
	}, nil

}
