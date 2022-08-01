package logic

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser() (resp *types.Response, err error) {
	//定义value 为string 类型的字典，用来存合约hash,合约编译器，文件名字
	var profile = make(map[string]interface{})

	reader, err := l.r.MultipartReader()

	//根据当前时间戳来创建文件夹，用来存放合约作者要上传的合约源文件
	pathFile := "./img/"
	if err != nil {
		return &types.Response{Message: "faild"},err
	}
	// 读取作者上传的文件以及ContractHash,CompilerVersion等数据，并保存在map中。
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			data, _ := ioutil.ReadAll(part)
			if part.FormName() == "UserName" {
				profile[part.FormName()] = string(data)
			} else if part.FormName() == "Address" {
				//创建以address为名的目录
				path := "./img/"+part.FormName()
				isExit := isDirExists(path)
				if isExit{
					pathFile = path
				}else {
					os.Mkdir(path, 0777)
				}
				profile[part.FormName()] = string(data)
			} else if part.FormName() == "Email" {
				profile[part.FormName()] = string(data)
			}else if part.FormName() == "Twitter" {
				profile[part.FormName()] = string(data)
			} else if part.FormName() == "Bio" {
				profile[part.FormName()] = string(data)
			}else if part.FormName() == "Signature" {
				profile[part.FormName()] = string(data)
			}
		} else {
			dst, err := os.OpenFile(pathFile+"/"+part.FileName(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
			if err != nil{
				return &types.Response{Message: "faild"}, err
			}
			defer func(dst *os.File) {
				err := dst.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
			}(dst)
			_, err = io.Copy(dst, part)
			if err != nil {
				return &types.Response{Message: "faild"},err
			}

		}
	}


	return
}



//目录是否存在
func isDirExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
