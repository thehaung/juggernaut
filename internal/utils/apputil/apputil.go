package apputil

import "github.com/thehaung/juggernaut/internal/logger"

func Recovery() {
	if err := recover(); err != nil {
		logger.GetLogger().Error(err)
	}
}
