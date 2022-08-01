package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/keys"
	"io"
	"io/ioutil"
	"magaOasis/lib/type/nullstring"
	"magaOasis/model/user"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"magaOasis/internal/svc"
	"magaOasis/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadUserLogic(ctx context.Context,r *http.Request, svcCtx *svc.ServiceContext) *UploadUserLogic {
	return &UploadUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		r:r,
		svcCtx: svcCtx,
	}
}

func (l *UploadUserLogic) UploadUser() (resp *types.Response, err error) {
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
		}else if part.FormName() == "Avatar"{
			pathFile := createDateDir("./")
			suffix := path.Ext(part.FileName())
			pf := pathFile+"/"+part.FormName()+suffix
			profile[part.FormName()] = pf
			dst, _ := os.OpenFile(pf,os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
			defer func(dst *os.File) {
				err := dst.Close()
				if err != nil {
					fmt.Printf("Closing file error: %v\n", err)
				}
			}(dst)
			_, err := io.Copy(dst, part)
			if err != nil {
				fmt.Printf("Copy error: %v\n", err)
				return nil, err
			}
			dst.Close()
		}else if part.FormName() == "Banner"{
			pathFile := createDateDir("./")
			suffix := path.Ext(part.FileName())
			pf := pathFile+"/"+part.FormName()+suffix
			profile[part.FormName()] = pf
			dst, _ := os.OpenFile(pf ,os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
			defer func(dst *os.File) {
				err := dst.Close()
				if err != nil {
					fmt.Printf("Closing file error: %v\n", err)
				}
			}(dst)
			_, err := io.Copy(dst, part)
			if err != nil {
				fmt.Printf("Copy error: %v\n", err)
				return nil ,err
			}

			dst.Close()
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
		banner := profile["Banner"].(string)
		avatar := profile["Avatar"].(string)
		bannerName := path.Base(banner)
		avatarName := path.Base(avatar)

		// 处理上传的文件
		//创建以address为名的目录
		pathDir := "./img/"+profile["Address"].(string)
		isExit := isDirExists(pathDir)
		if !isExit{    // 如果不存在，则创建新目录
			os.Mkdir(pathDir, 0777)
		}

		bannerOldFileName:=banner
		avatarOldFileName:=avatar

		bannerNewFileName:=pathDir+"/"+bannerName
		avatarNewFileName:=pathDir+"/"+avatarName
		err := os.Rename(bannerOldFileName,bannerNewFileName)
		if err!=nil{
			fmt.Println("重命名失败",err)
		}
		err = os.Rename(avatarOldFileName,avatarNewFileName)
		if err!=nil{
			fmt.Println("重命名失败",err)
		}
		pD := path.Dir(banner)
		err = os.Remove(pD)
        if err !=nil{
			fmt.Println("remove dir failed",err)
			return nil, err
		}

		//save img name
		bannerFullname := "/"+profile["Address"].(string) +"/"+ bannerName
		avatarFullname := "/"+ profile["Address"].(string) + "/"+avatarName

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


func createDateDir(basepath string) string {
	folderName := time.Now().Format("20060102150405")

	fmt.Println("Create folder: %v", folderName)
	folderPath := filepath.Join(basepath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0777)
		if err != nil {
			fmt.Println("Create dir error: %v", err)
		}
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			fmt.Println("Chmod error: %v", err)
		}
	}
	return folderPath
}