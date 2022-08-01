package logic

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"magaOasis/internal/svc"
	"magaOasis/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageLogic {
	return &GetImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetImageLogic) GetImage(req *types.FilePath) (resp *types.FileResponse, err error) {

	//iot,err := ioutil.ReadFile("D:\\Neo\\GoProject\\megaOasis\\img\\wangmingting\\avatar.nef") //读取文件   img/wangmingting/avatar.nef

	//打开文件
	file, err := os.Open("D:\\Neo\\GoProject\\megaOasis\\img\\wangmingting\\avatar.nef") //只是用来读的时候，用os.Open。相对路径，针对于同目录下。
	if err != nil{
		fmt.Printf("打开文件失败,err:%v\n",err)
		return
	}
	defer file.Close() //关闭文件,为了避免文件泄露和忘记写关闭文件


	 content, err := ioutil.ReadAll(file)

    fmt.Println(string(content))



	return &types.FileResponse{Code:content},nil
}
