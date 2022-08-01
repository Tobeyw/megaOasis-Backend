package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"io/ioutil"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"net/http"
	"os"
)

type TestUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewTestUserLogic(r *http.Request, svcCtx *svc.ServiceContext) *TestUserLogic {
	return &TestUserLogic{
		Logger: logx.WithContext(r.Context()),
		r:    r,
		svcCtx: svcCtx,
	}
}

func (l *TestUserLogic) TestUser() (resp *types.Response, err error) {
	var profile = make(map[string]interface{})
	reader, err := l.r.MultipartReader()
	if err != nil{
		return &types.Response{Message: "file upload failed,err:"},err
	}
	if err != nil {
		return &types.Response{Message: "failed"},err
	}
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
		} else if part.FormName() == "Avatar"{
			profile[part.FormName()] = part
		}else if part.FormName()== "Banner"{
			profile[part.FormName()] = part
		}
	}

	//处理验签
	signatureData := profile["Signature"].(string)
	var signature types.SignatureData
	if err1 := json.Unmarshal([]byte(signatureData), &signature); err1 != nil {
		return &types.Response{Message: "signature convert failed"},err
	}

	// 处理上传的文件
	//创建以address为名的目录
	pathDir := "./img/"+profile["Address"].(string)
	isExit := isDirExists(pathDir)
	if !isExit{    // 如果不存在，则创建新目录
		os.Mkdir(pathDir, 0777)
	}





	return
}

