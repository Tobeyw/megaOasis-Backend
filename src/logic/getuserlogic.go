package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	neo "magaOasis/common/mongo"
	"magaOasis/src/config"
	"magaOasis/src/svc"
	"magaOasis/src/types"
	"os"

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

	fmt.Println(req.Address)
	res, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, req.Address)
	fmt.Println(err)
	if err != nil && err != sqlc.ErrNotFound {

		return nil, err
	}
	fmt.Println(err)
	if err == sqlc.ErrNotFound {
		return nil, nil
	}

	cd, dbonline := intializeMongoOnlineClient(l.svcCtx.Config, context.TODO())
	me := neo.T{
		Db_online: dbonline,
		C_online:  cd,
	}
	isValid, err := me.IsOwnerByNNS(res.NNS.String, res.Address)
	if err != nil {
		return nil, err
	}
	if !isValid {
		res.NNS.String = ""
	}
	return &types.UserResp{
		UserName: res.Username.String,
		Address:  res.Address,
		NNS:      res.NNS.String,
		Email:    res.Email.String,
		Twitter:  res.Twitter.String,
		Avatar:   res.Avatar.String,
		Bio:      res.Bio.String,
		Banner:   res.Banner.String,
	}, nil

}

func intializeMongoOnlineClient(cfg config.Config, ctx context.Context) (*mongo.Client, string) {
	rt := os.ExpandEnv("${RUNTIME}")
	//默认main
	clientOptions := options.Client().ApplyURI(cfg.MongoDBMain)
	dbOnline := cfg.DBMain
	if rt == "test" {
		clientOptions = options.Client().ApplyURI(cfg.MongoDBTest)
		dbOnline = cfg.DBTest
	} else if rt == "main" {
		clientOptions = options.Client().ApplyURI(cfg.MongoDBMain)
		dbOnline = cfg.DBMain
	} else if rt == "dev" {
		clientOptions = options.Client().ApplyURI(cfg.MongoDBDev)
		dbOnline = cfg.DBDev
	}

	clientOptions.SetMaxPoolSize(20)
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	return co, dbOnline
}
