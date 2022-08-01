package logic

import (
	"context"
	"database/sql"
	"magaOasis/lib/type/nullstring"
	"magaOasis/model/user"
	"time"

	"magaOasis/internal/svc"
	"magaOasis/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.UserReq) (*types.Response,error) {

	_, err := l.svcCtx.UserModel.Insert(l.ctx,&user.User{
		Username: sql.NullString{req.UserName,nullstring.IsNull(req.UserName)},
		Address: req.Address,
		Bio: sql.NullString{req.Bio,nullstring.IsNull(req.Bio)},
		Email:sql.NullString{req.Email,nullstring.IsNull(req.Email)} ,
		Twitter: sql.NullString{req.Twitter,nullstring.IsNull(req.Twitter)},
		Avatar: sql.NullString{req.Avatar,nullstring.IsNull(req.Avatar)},
		Banner: sql.NullString{req.Banner,nullstring.IsNull(req.Avatar)},
		Timestamp: time.Now().UnixMilli(),
	})

	if err != nil {
		return nil, err
	}

	return &types.Response{Message: "success"}, nil
}
