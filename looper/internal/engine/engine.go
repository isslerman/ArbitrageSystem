package engine

import (
	"fmt"
	"looper/config"
	"looper/infra/ntfy"
	"looper/internal/marketdata"
	"looper/internal/repository"
	"looper/internal/strategy"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Our app confiiguration.
// Errors and Infos for debug: All the erros and info msgs are saved to a DB.
type Engine struct {
	cfg        *config.Config // Configuration for the app
	Repo       repository.Repository
	Notify     *ntfy.Ntfy // Notify sends notifications to mobile
	MarketData *marketdata.MarketData
	Strategy   strategy.Strategy
	Logger     *zap.SugaredLogger
	// ac       *ArbitrageControld
}

func New(cfg *config.Config) *Engine {
	e := &Engine{
		cfg: cfg,
	}
	e.Logger = e.initLogger()
	return e
}

func (e *Engine) initLogger() *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")

	// 💡 Enable color on the log level
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}

	return zapLogger.Sugar()
}

func (e *Engine) GracefullyShutdown() {
	e.Logger.Info("Shutting down server...")
	e.Logger.Info("Closing all open orders.")

	e.Repo.LogInfo("ALERT: Shutting down server...")
	e.Repo.LogInfo("ALERT: Closing all open orders.")

	if err := e.Strategy.StopAllOrders(); err != nil {
		e.Logger.Errorf("Error canceling all orders, check the exchange for open orders.", "err", err)
	} else {
		e.Logger.Info("All orders cancelled successfully.")
	}

	e.Repo.Stop()
	e.Logger.Info("Shutdown complete")
}

// setupNotifyService send a msg to ntfy mobile app
func (e *Engine) setupNotifyService() {
	// TEMPORARILY DISABLED - REMOVE THE TWO COMMENTS TO ENABLE

	// app.Notify = ntfy.NewNtfy()
	// app.Notify.SendMsg("Server is up - NTFY Test", "Server is up - LoggerErr Test", false)
}
