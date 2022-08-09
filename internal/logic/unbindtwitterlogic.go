package logic

import (
	"context"
	"database/sql"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"magaOasis/lib/type/nullstring"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindTwitterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindTwitterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindTwitterLogic {
	return &UnbindTwitterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindTwitterLogic) UnbindTwitter(req *types.UnbindTwitter) (resp *types.Response, err error) {
	verifyResult, err := isVerify(req.Signature)
	if err != nil {
		return &types.Response{Code: 32001, Message: "signature verify failed"}, err
	}

	if verifyResult {
		res, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, req.Address)
		if err != nil {
			return &types.Response{Code: 32001, Message: "signature verify failed"}, err
		}

		//fmt.Println(req.Twitter,res.Twitter.String,req.Twitter==res.Twitter.String)
		if req.Twitter == res.Twitter.String {
			res.Twitter = sql.NullString{"", nullstring.IsNull("")}
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
