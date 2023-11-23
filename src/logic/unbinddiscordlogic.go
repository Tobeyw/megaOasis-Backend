package logic

import (
	"context"
	"database/sql"
	"magaOasis/lib/type/nullstring"
	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindDiscordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindDiscordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindDiscordLogic {
	return &UnbindDiscordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindDiscordLogic) UnbindDiscord(req *types.UnbindDiscord) (resp *types.Response, err error) {
	verifyResult, err := isVerify(req.Signature)
	if err != nil {
		return &types.Response{Code: 32001, Message: "signature verify failed"}, err
	}

	if verifyResult {
		res, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, req.Address)
		if err != nil {
			return &types.Response{Code: 32001, Message: "signature verify failed"}, err
		}

		if req.Discord == res.Discord.String {
			res.Discord = sql.NullString{"", nullstring.IsNull("")}
			err := l.svcCtx.UserModel.Update(l.ctx, res)
			if err != nil {
				return nil, err
			}
		} else {
			return &types.Response{Code: 32001, Message: "failed"}, nil
		}
	} else {
		return &types.Response{Code: 32001, Message: "signature verify false"}, nil
	}

	return &types.Response{Code: 200, Message: "success"}, nil
}
