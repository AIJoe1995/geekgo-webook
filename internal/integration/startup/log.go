package startup

import "geekgo-webook/pkg/logger"

func InitLog() logger.LoggerV1 {
	return &logger.NopLogger{}
}
