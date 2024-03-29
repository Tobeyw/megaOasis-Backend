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
	neo "magaOasis/common/mongo"
	"magaOasis/lib/type/nullstring"
	"magaOasis/model/user"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"magaOasis/src/svc"
	"magaOasis/src/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadUserLogic(ctx context.Context, r *http.Request, svcCtx *svc.ServiceContext) *UploadUserLogic {
	return &UploadUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *UploadUserLogic) UploadUser() (resp *types.Response, err error) {
	//读取上传的数据
	var profile = make(map[string]interface{})
	reader, err := l.r.MultipartReader()
	if err != nil {
		return &types.Response{Code: 32001, Message: "file upload failed,err:"}, err
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		b := part.FileName()
		fmt.Println(b)
		a := part.FormName()
		fmt.Println(a)
		if part.FileName() == "" {
			a := part.FormName()
			fmt.Println(a)
			data, _ := ioutil.ReadAll(part)
			profile[part.FormName()] = string(data)
			if part.FormName() == "Avatar" {
				profile["flag"+part.FormName()] = false
			} else if part.FormName() == "Banner" {
				profile["flag"+part.FormName()] = false
			}
			fmt.Println("===1:", profile["flag"+part.FormName()])
		} else {
			if part.FormName() == "Avatar" {
				pathFile := createDateDir("./")
				suffix := path.Ext(part.FileName())
				pf := pathFile + "/" + part.FormName() + suffix
				profile[part.FormName()] = pf
				profile["flag"+part.FormName()] = true
				fmt.Println("====avatar:", profile["flag"+part.FormName()])
				dst, _ := os.OpenFile(pf, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
				defer func(dst *os.File) {
					err := dst.Close()
					if err != nil {
						fmt.Printf("Closing file error: %v\n", err)
					}
				}(dst)
				_, err := io.Copy(dst, part)
				if err != nil {
					fmt.Printf("Copy error: %v\n", err)
					return &types.Response{Code: 32001, Message: "Copy error:"}, err
				}
				dst.Close()
			} else if part.FormName() == "Banner" {
				pathFile := createDateDir("./")
				suffix := path.Ext(part.FileName())
				pf := pathFile + "/" + part.FormName() + suffix
				profile[part.FormName()] = pf
				profile["flag"+part.FormName()] = true
				fmt.Println("====banner:", profile["flag"+part.FormName()])
				dst, _ := os.OpenFile(pf, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
				defer func(dst *os.File) {
					err := dst.Close()
					if err != nil {
						fmt.Printf("Closing file error: %v\n", err)
					}
				}(dst)
				_, err := io.Copy(dst, part)
				if err != nil {
					fmt.Printf("Copy error: %v\n", err)
					return &types.Response{Code: 32001, Message: "Copy error:"}, err
				}

				dst.Close()
			}
		}
	}

	//判断为空的情况
	UserName := ""
	if profile["UserName"] != nil {
		UserName = profile["UserName"].(string)
	}
	//
	Bio := ""
	if profile["Bio"] != nil {
		Bio = profile["Bio"].(string)
	}
	Address := ""
	if profile["Address"] != nil {
		Address = profile["Address"].(string)
	}
	Twitter := ""
	if profile["Twitter"] != nil {
		Twitter = profile["Twitter"].(string)
	}
	Discord := ""
	if profile["Discord"] != nil {
		Discord = profile["Discord"].(string)
	}
	Email := ""
	if profile["Email"] != nil {
		Email = profile["Email"].(string)
	}
	NNS := ""
	if profile["NNS"] != nil {
		NNS = profile["NNS"].(string)
	}
	getUser, err := l.svcCtx.UserModel.FindOneByAddress(l.ctx, Address)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return &types.Response{Code: 32001, Message: "findByAddress err:"}, err
	}
	//处理username 重复的问题
	if getUser == nil && UserName != "" { // add
		getUserByName, _ := l.svcCtx.UserModel.FindOneByUserName(l.ctx, UserName)
		if getUserByName != nil {
			return &types.Response{Code: 32002, Message: "name already exists"}, nil
		}

	}
	//校验nns的owner 是否为address & 检查nns 是否已经被设置，保证唯一
	if NNS != "" {
		//取消owner校验  （nns在转出，被设置之前，依旧可以展示在之前owner的propfile中）
		cd, dbonline := intializeMongoOnlineClient(l.svcCtx.Config, context.TODO())
		me := neo.T{
			Db_online: dbonline,
			C_online:  cd,
		}
		isValid, err := me.IsOwnerByNNS(NNS, Address)
		if err != nil {
			return &types.Response{Code: 32001, Message: "nns invalid parameter"}, err
		}
		if isValid {
			//唯一性校验
			preUser, err := l.svcCtx.UserModel.FindOneByNNS(l.ctx, NNS)
			if err != nil {
				return &types.Response{Code: 32002, Message: "query user by nns err"}, err
			}

			if preUser != nil { //nns 的owner发生改变  并被新的owner设置，之前的owne 的profile 旧的owner的nns需要重置
				preUser.NNS = sql.NullString{"", nullstring.IsNull("")}
				err := l.svcCtx.UserModel.Update(l.ctx, preUser)
				if err != nil {
					return &types.Response{Code: 32001, Message: "update failed"}, err
				}
			}
		} else {
			NNS = ""
		}

	}

	//处理验签数据
	signatureData := profile["Signature"].(string)
	var signature types.SignatureData
	if err1 := json.Unmarshal([]byte(signatureData), &signature); err1 != nil {
		return &types.Response{Code: 32001, Message: "signature convert failed"}, err
	}

	verifyResult, err := isVerify(signature)
	if err != nil {
		return &types.Response{Code: 32001, Message: "signature verify failed"}, err
	}

	if verifyResult {
		removeDir := ""
		bannerFullname := ""
		avatarFullname := ""
		// 处理上传的文件
		//创建以address为名的目录
		pathDir := "./img/" + profile["Address"].(string)
		isExit := isDirExists(pathDir)
		if !isExit { // 如果不存在，则创建新目录
			os.Mkdir(pathDir, 0777)
		}
		if profile["flagBanner"].(bool) {
			banner := profile["Banner"].(string)
			removeDir = banner
			bannerName := path.Base(banner)
			bannerOldFileName := banner
			bannerNewFileName := pathDir + "/" + bannerName
			var cmd *exec.Cmd
			cmd = exec.Command("mv", bannerOldFileName, bannerNewFileName)
			_, err := cmd.Output()
			if err != nil {
				return &types.Response{Code: 32001, Message: "rename failed:"}, err
			}
			bannerFullname = "/" + profile["Address"].(string) + "/" + bannerName
		} else {
			bannerFullname = profile["Banner"].(string)
		}

		if profile["flagAvatar"].(bool) {
			avatar := profile["Avatar"].(string)
			removeDir = avatar
			avatarName := path.Base(avatar)
			avatarOldFileName := avatar
			avatarNewFileName := pathDir + "/" + avatarName
			var cmd *exec.Cmd
			cmd = exec.Command("mv", avatarOldFileName, avatarNewFileName)
			_, err := cmd.Output()
			if err != nil {
				return &types.Response{Code: 32001, Message: "rename failed:"}, err
			}
			avatarFullname = "/" + profile["Address"].(string) + "/" + avatarName
		} else {
			avatarFullname = profile["Avatar"].(string)
		}

		if removeDir != "" {
			pD := path.Dir(removeDir)
			err = os.Remove(pD)
			if err != nil {
				fmt.Println("remove dir failed", err)
				return &types.Response{Code: 32001, Message: "remove dir failed"}, err
			}
		}

		EmailCreate := int64(0)
		//将数据存入数据库
		if Email != "" {
			EmailCreate = time.Now().UnixMilli()
		}

		if getUser == nil { //create
			_, err := l.svcCtx.UserModel.Insert(l.ctx, &user.User{
				Username:      sql.NullString{UserName, nullstring.IsNull(UserName)},
				Address:       Address,
				NNS:           sql.NullString{NNS, nullstring.IsNull(NNS)},
				Bio:           sql.NullString{Bio, nullstring.IsNull(Bio)},
				Email:         sql.NullString{Email, nullstring.IsNull(Email)},
				Twitter:       sql.NullString{Twitter, nullstring.IsNull(Twitter)},
				Discord:       sql.NullString{Discord, nullstring.IsNull(Discord)},
				Avatar:        sql.NullString{avatarFullname, nullstring.IsNull(avatarFullname)},
				Banner:        sql.NullString{bannerFullname, nullstring.IsNull(bannerFullname)},
				Timestamp:     time.Now().UnixMilli(),
				TwitterCreate: sql.NullInt64{int64(0), nullstring.IsZero(int64(0))},
				EmailCreate:   sql.NullInt64{EmailCreate, nullstring.IsZero(EmailCreate)},
			})

			if err != nil {

				return &types.Response{Code: 32001, Message: "insert failed"}, err
			}

			return &types.Response{Code: 200, Message: "success"}, nil
		} else { //update

			err := l.svcCtx.UserModel.Update(l.ctx, &user.User{
				Id:            getUser.Id,
				Username:      sql.NullString{UserName, nullstring.IsNull(UserName)},
				Address:       Address,
				NNS:           sql.NullString{NNS, nullstring.IsNull(NNS)},
				Bio:           sql.NullString{Bio, nullstring.IsNull(Bio)},
				Email:         sql.NullString{Email, nullstring.IsNull(Email)},
				Twitter:       sql.NullString{Twitter, nullstring.IsNull(Twitter)},
				Discord:       sql.NullString{Discord, nullstring.IsNull(Discord)},
				Avatar:        sql.NullString{avatarFullname, nullstring.IsNull(avatarFullname)},
				Banner:        sql.NullString{bannerFullname, nullstring.IsNull(bannerFullname)},
				Timestamp:     getUser.Timestamp,
				EmailCreate:   sql.NullInt64{EmailCreate, nullstring.IsZero(EmailCreate)},
				TwitterCreate: getUser.TwitterCreate,
			})

			if err != nil {
				return &types.Response{Code: 32001, Message: "update failed"}, err
			}

			return &types.Response{Code: 200, Message: "success"}, nil

		}
	} else {
		return &types.Response{Code: 32001, Message: "faild"}, nil
	}

	return &types.Response{}, err
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

// 目录是否存在
func isDirExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// 验签
func isVerify(signature types.SignatureData) (bool, error) {

	//验签
	publicKey := signature.PublicKey
	pk, err := crypto.NewECPointFromString(publicKey)
	if err != nil {
		return false, err
	}
	data := helper.HexToBytes(signature.Data)
	parameterHexString := helper.BytesToHex([]byte(signature.Salt + signature.Message))

	varint := helper.VarIntFromInt(len(parameterHexString) / 2)
	lengthHex := helper.BytesToHex(varint.Bytes())
	concatenatedString := lengthHex + parameterHexString
	serializedTransaction := "010001f0" + concatenatedString + "0000"

	result := keys.VerifySignature(helper.HexToBytes(serializedTransaction), data, pk)
	return result, nil
}
