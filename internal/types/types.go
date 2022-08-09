// Code generated by goctl. DO NOT EDIT.
package types

type Address struct {
	Address string `form:"address"`
}

type UserReq struct {
	UserName  string        `json:"username"`
	Address   string        `json:"address"`
	Email     string        `json:"email"`
	Twitter   string        `json:"twitter"`
	Avatar    string        `json:"avatar"`
	Bio       string        `json:"bio"`
	Banner    string        `json:"banner"`
	Signature SignatureData `json:"signature"`
}

type UserResp struct {
	UserName      string `json:"username"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	Twitter       string `json:"twitter"`
	Avatar        string `json:"avatar"`
	Bio           string `json:"bio"`
	Banner        string `json:"banner"`
	Timestamp     int64  `json:"timestamp"`
	TtitterUpdate int64  `json:"twitterupdate"`
	EmailUpdate   int64  `json:"emailupdate"`
}

type SignatureData struct {
	Message   string `json:"message"`
	PublicKey string `json:"publicKey"`
	Salt      string `json:"salt"`
	Data      string `json:"data"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type FilePath struct {
	Path string `form:"path"`
}

type FileResponse struct {
	Code []byte `json:"code"`
}

type LoginTwitterResponse struct {
	Url string `json:"url"`
}

type LoginTwitterParam struct {
	Address     string `json:"address" form:"address"`
	CallbackUrl string `json:"callbackUrl" form:"callbackUrl"`
}

type CallbackTwitterParam struct {
	State string `form:"state"`
	Code  string `form:"code"`
}

type UnbindTwitter struct {
	Address   string        `json:"address"`
	Twitter   string        `json:"twitter"`
	Signature SignatureData `json:"signature"`
}

type TwitterAccessToken struct {
	Address   string        `json:"address"`
	Code      string        `json:"code"`
	Signature SignatureData `json:"signature"`
}
