package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ctx removed for this app
// type ctxKey struct{}

var once sync.Once

var logger *zap.Logger

// Get initializes a zap.Logger instance if it has not been initialized
// already and returns the same instance for subsequent calls.
// It receives an ID that is used to create different log filenames
func Get(id string) *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/%s.log", id),
			MaxSize:    5,     // max size in MB
			MaxBackups: 10,    // max number of log files to retain
			MaxAge:     1,     // number max of days to retain old logs
			Compress:   false, // compress using gzip
		})

		level := zap.InfoLevel
		levelEnv := os.Getenv("LOG_LEVEL")
		if levelEnv != "" {
			levelFromEnv, err := zapcore.ParseLevel(levelEnv)
			if err != nil {
				log.Println(
					fmt.Errorf("invalid level, defaulting to INFO: %w", err),
				)
			}

			level = levelFromEnv
		}

		logLevel := zap.NewAtomicLevelAt(level)

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		// log to multiple destinations (console and file)
		// extra fields are added to the JSON output alone
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel),
			zapcore.NewCore(fileEncoder, file, logLevel).
				With(
					[]zapcore.Field{
						zap.String("git_revision", gitRevision),
						zap.String("go_version", buildInfo.GoVersion),
					},
				),
		)

		logger = zap.New(core)
	})

	return logger
}

// // FromCtx returns the Logger associated with the ctx. If no logger
// // is associated, the default logger is returned, unless it is nil
// // in which case a disabled logger is returned.
// func FromCtx(ctx context.Context) *zap.Logger {
// 	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
// 		return l
// 	} else if l := logger; l != nil {
// 		return l
// 	}

// 	return zap.NewNop()
// }

// // WithCtx returns a copy of ctx with the Logger attached.
// func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
// 	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
// 		if lp == l {
// 			// Do not store same logger.
// 			return ctx
// 		}
// 	}

// 	return context.WithValue(ctx, ctxKey{}, l)
// }
