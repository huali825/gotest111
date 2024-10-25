package startup

import (
	"goworkwebook/webook003/internal/service/oauth2/wechat"
	"goworkwebook/webook003/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
