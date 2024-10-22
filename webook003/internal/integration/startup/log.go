package startup

import "goworkwebook/webook003/pkg/logger"

func InitLogger() logger.LoggerV1 {
	return logger.NewNopLogger()
}
