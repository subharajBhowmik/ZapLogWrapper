package NISZapLogWrapper

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func InitLogger(outputPaths []string) *zap.Logger {

	logConfig := zap.NewProductionConfig()

	// CHANGE TO LOG FILE PATH ONLY
	//logConfig.OutputPaths = []string{"/Users/subharajbhowmik/go/src/NISAuthenticationService/dummyLog.log", "stderr"}
	logConfig.OutputPaths = outputPaths

	logConfig.EncoderConfig.FunctionKey = "func"
	logConfig.EncoderConfig.TimeKey = "time"

	logConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700")) // Change to time.Stamp for: Jun 24 11:37:42
	}

	logConfig.DisableStacktrace = true
	logger, _ := logConfig.Build()

	return logger
}
