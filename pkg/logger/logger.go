package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Module = fx.Options(
	fx.Provide(ProvideLogger),
)

func ProvideLogger() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: "log/app.log",
		MaxSize:  100,
		MaxAge:   60,
	})

	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "datetime"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.DebugLevel,
	)

	return zap.New(core)
}
