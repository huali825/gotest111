package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/pkg/logger"
	"net/http"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	appID     string
	appSecret string
	client    *http.Client
	l         logger.LoggerV1
}

func NewService(appID string, appSecret string, l logger.LoggerV1) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
		l:         l,
	}
}

// AuthURL 构建扫码登录url
func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}

// VerifyCode 函数用于验证微信的授权码
func (s *service) VerifyCode(ctx context.Context,
	code string) (domain.WechatInfo, error) {
	// 构造获取 access_token 的 URL
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	// 创建 GET 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	// 发送请求
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	var res Result
	// 解析响应
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		// 转 JSON 为结构体出错
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("调用微信接口失败 errcode %d, errmsg %s", res.ErrCode, res.ErrMsg)
	}
	// 返回微信信息
	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
