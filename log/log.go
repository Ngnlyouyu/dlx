package log

import "go.uber.org/zap"

var Sugar *zap.SugaredLogger

func InitSugar() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("init sugared logger failed: " + err.Error())
	}
	Sugar = logger.Sugar()
}
