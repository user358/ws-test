package log

import "go.uber.org/zap"

var log *zap.SugaredLogger

func Init() func() {
	logger, _ := zap.NewProduction(zap.AddCallerSkip(1))

	log = logger.Sugar()

	return func() {
		_ = logger.Sync()
	}
}

func Infow(msg string, keysAndValues ...interface{}) {
	log.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	log.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	log.Fatalw(msg, keysAndValues...)
}
