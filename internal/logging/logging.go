package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() *zap.Logger {
	// info level enabler
	level := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.DebugLevel
	})

	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stdoutSyncer,
			level,
		),
	)

	return zap.New(core)
}
