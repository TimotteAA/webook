package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"webook/internal/domain"
)

// 扫码url来自于微信文档
const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

// 扫码成功后的回调地址
var redirectURL = url.PathEscape("https://xxx.com/oauth2/wechat/callback")

type WeChatService interface {
	AuthURL(state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WeChatResult, error)
}

type wechatService struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewWeChatService(appId string, appSecret string) WeChatService {
	return &wechatService{appId: appId, appSecret: appSecret, client: http.DefaultClient}
}

func (w *wechatService) AuthURL(state string) (string, error) {

	return fmt.Sprintf(authURLPattern, w.appId, redirectURL, state), nil
}

func (w *wechatService) VerifyCode(ctx context.Context, code string) (domain.WeChatResult, error) {
	//	调微信接口判断
	const baseUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
	queryParams := url.Values{}
	queryParams.Set("appid", w.appId)
	queryParams.Set("secret", w.appSecret)
	queryParams.Set("code", code)
	queryParams.Set("grant_type", "authorization_code")
	url := baseUrl + "?" + queryParams.Encode()
	// get请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.WeChatResult{}, err
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return domain.WeChatResult{}, err
	}
	defer resp.Body.Close()
	// 解码body
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return domain.WeChatResult{}, err
	}
	// 错误码
	if result.ErrCode != 0 {
		return domain.WeChatResult{}, errors.New("换取微信access_token失败")
	}
	return domain.WeChatResult{UnionId: result.UnionId, OpenId: result.OpenId}, nil
}

// 微信拿token返回的字段
type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errMsg"`

	Scope string `json:"scope"`

	AccessToken string `json:"access_token"`
	// 过期时间
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	// 此处只关心下面的两个id，来标识用户
	// 应用下的id
	OpenId string `json:"openid"`
	// 公司下的id
	UnionId string `json:"unionid"`
}
