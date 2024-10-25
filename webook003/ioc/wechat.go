package ioc

import (
	"goworkwebook/webook003/internal/service/oauth2/wechat"
	"goworkwebook/webook003/pkg/logger"
)

// InitWechatService 这里web不直接用service的NewService , 是因为要从环境变量拿参数
func InitWechatService(l logger.LoggerV1) wechat.Service {
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("找不到环境变量 WECHAT_APP_ID")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("找不到环境变量 WECHAT_APP_SECRET")
	//}
	appID := "11"
	appSecret := "22"
	return wechat.NewService(appID, appSecret, l)
}
