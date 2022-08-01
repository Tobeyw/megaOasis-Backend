package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/keys"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"io/ioutil"
	"magaOasis/internal/svc"
	"magaOasis/internal/types"
	"magaOasis/lib/type/nullstring"
	"magaOasis/model/user"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"
)

type EditUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewEditUserLogic(r *http.Request,svcCtx *svc.ServiceContext) *EditUserLogic {
	return &EditUserLogic{
		Logger: logx.WithContext(r.Context()),
		r:    r,
		svcCtx: svcCtx,
	}
}

func (l *EditUserLogic) EditUser() (resp *types.Response, err error) {
	//读取上传的数据
	var profile = make(map[string]interface{})
	reader, err := l.r.MultipartReader()
	if err != nil{
		return &types.Response{Message: "file upload failed,err:"},err
	}
	if err != nil {
		return &types.Response{Message: "faild"},err
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
			fmt.Println("==============",part)
		}else if part.FormName()== "Banner"{
			profile[part.FormName()] = part
		}
	}
	//处理验签数据
	signatureData := profile["Signature"].(string)
	var signature types.SignatureData
	if err1 := json.Unmarshal([]byte(signatureData), &signature); err1 != nil {
		return &types.Response{Message: "signature convert failed"},err
	}

	//验签
	publicKey := signature.PublicKey
	pk,err := crypto.NewECPointFromString(publicKey)
	if err != nil{
		return nil, err
	}
	data := helper.HexToBytes(signature.Data)
	parameterHexString :=helper.BytesToHex([]byte( signature.Salt +signature.Message))

	varint := helper.VarIntFromInt(len(parameterHexString)/2)
	lengthHex := helper.BytesToHex(varint.Bytes())
	concatenatedString := lengthHex + parameterHexString
	serializedTransaction := "010001f0" + concatenatedString + "0000";

	isVerify :=	keys.VerifySignature(helper.HexToBytes(serializedTransaction),data,pk)

	if isVerify {
		// 处理上传的文件
		//创建以address为名的目录
		pathDir := "./img/"+profile["Address"].(string)
		isExit := isDirExists(pathDir)
		if !isExit{    // 如果不存在，则创建新目录
			os.Mkdir(pathDir, 0777)
		}

		banner := profile["Banner"].(*multipart.Part)
		avatar := profile["Avatar"].(*multipart.Part)
		bannerSuffix :=  path.Ext(banner.FileName())
		avatarSuffix := path.Ext(avatar.FileName())

		bannerFullname := pathDir + "/banner" + bannerSuffix
		avatarFullname := pathDir + "/avatar" + avatarSuffix



		UserName := profile["UserName"].(string)
		Bio := profile["Bio"].(string)
		Address := profile["Address"].(string)
		Twitter := profile["Twitter"].(string)
		Email := profile["Email"].(string)
		//将数据存入数据库
		getUser, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx,Address)
		if err != nil && err.Error() != "sql: no rows in result set"{
			return nil, err
		}

		if getUser == nil{             //create
			_, err := l.svcCtx.UserModel.Insert(l.ctx,&user.User{
				Username: sql.NullString{UserName,nullstring.IsNull(UserName)},
				Address: Address,
				Bio: sql.NullString{Bio,nullstring.IsNull(Bio)},
				Email:sql.NullString{Email,nullstring.IsNull(Email)} ,
				Twitter: sql.NullString{Twitter,nullstring.IsNull(Twitter)},
				Avatar: sql.NullString{avatarFullname,nullstring.IsNull(avatarFullname)},
				Banner: sql.NullString{bannerFullname,nullstring.IsNull(bannerFullname)},
				Timestamp: time.Now().UnixMilli(),
			})

			if err != nil {
				return nil, err
			}

			return &types.Response{Message: "success"}, nil
		}else{                   //update

			err := l.svcCtx.UserModel.Update(l.ctx,&user.User{
				Id: getUser.Id,
				Username: sql.NullString{UserName,nullstring.IsNull(UserName)},
				Address: Address,
				Bio: sql.NullString{Bio,nullstring.IsNull(Bio)},
				Email:sql.NullString{Email,nullstring.IsNull(Email)} ,
				Twitter: sql.NullString{Twitter,nullstring.IsNull(Twitter)},
				Avatar: sql.NullString{avatarFullname,nullstring.IsNull(avatarFullname)},
				Banner: sql.NullString{bannerFullname,nullstring.IsNull(bannerFullname)},
				Timestamp: getUser.Timestamp,
			})

			if err != nil {
				return nil, err
			}




			return &types.Response{Message: "success"}, nil

		}
	}else {
		return  &types.Response{Message: "faild"}, nil
	}



	return &types.Response{},err
}


