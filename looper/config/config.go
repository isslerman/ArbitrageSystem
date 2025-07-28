package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		GrpcPort int `mapstructure:"grpc_port"`
	} `mapstructure:"server"`

	Symbols struct {
		Base  string `mapstructure:"base"`
		Quote string `mapstructure:"quote"`
	} `mapstructure:"symbols"`

	Trade struct {
		Amount               float64 `mapstructure:"amount"`
		ThresholdRemakeOrder float64 `mapstructure:"threshold_remake_order"`
		Profit               float64 `mapstructure:"profit"`
		AskSide              bool    `mapstructure:"ask_side"`
		BidSide              bool    `mapstructure:"bid_side"`
	} `mapstructure:"trade"`

	Fees struct {
		Maker float64 `mapstructure:"maker"`
		Taker float64 `mapstructure:"taker"`
	} `mapstructure:"fees"`

	Websocket struct {
		WssURL string `mapstructure:"wss_url"`
	} `mapstructure:"websocket"`

	Database struct {
		DSN                    string `mapstructure:"dsn"`
		Host                   string `mapstructure:"host"`
		Port                   int    `mapstructure:"port"`
		User                   string `mapstructure:"user"`
		Password               string `mapstructure:"password"`
		Name                   string `mapstructure:"name"`
		SSLMode                string `mapstructure:"ssl_mode"`
		Timezone               string `mapstructure:"timezone"`
		ConnectTimeoutSeconds  int    `mapstructure:"connect_timeout_seconds"`
		MaxOpenConns           int    `mapstructure:"pool.max_open_conns"`
		MaxIdleConns           int    `mapstructure:"pool.max_idle_conns"`
		ConnMaxLifetimeMinutes int    `mapstructure:"pool.conn_max_lifetime_minutes"`
	} `mapstructure:"database"`
}

// Load YML config file - config/config.yml
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")

	// Defaults
	v.SetDefault("server.grpc_port", 50001)
	v.SetDefault("symbols.base", "USDT")
	v.SetDefault("symbols.quote", "BRL")
	v.SetDefault("trade.threshold_remake_order", 0.15)
	v.SetDefault("trade.profit", 1.00)
	v.SetDefault("trade.amount", 0.10)
	v.SetDefault("trade.ask_side", true)
	v.SetDefault("trade.bid_side", false)
	v.SetDefault("fees.fee_maker", 0.00450)
	v.SetDefault("fees.fee_taker", 0.00075)
	v.SetDefault("websocket.wss_url", "wss://stream.binance.com:9443/ws")

	// Database defaults, safe to override in file or env
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.timezone", "UTC")
	v.SetDefault("database.connect_timeout_seconds", 5)
	v.SetDefault("database.pool.max_open_conns", 20)
	v.SetDefault("database.pool.max_idle_conns", 10)
	v.SetDefault("database.pool.conn_max_lifetime_minutes", 30)

	// Env overrides for any key like SERVER_GRPC_PORT
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Optional direct env binds for secrets
	_ = v.BindEnv("database.password")
	_ = v.BindEnv("database.user")
	_ = v.BindEnv("database.host")
	_ = v.BindEnv("database.port")
	_ = v.BindEnv("database.name")
	_ = v.BindEnv("database.dsn")

	// Read file
	if err := v.ReadInConfig(); err != nil {
		// If you want to run without a file, handle not found gracefully
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	// Unmarshal everything at once
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	slog.Info("Loaded config", "dsn", cfg.Database.DSN, "user", cfg.Database.User)

	// Only build DSN if it's not set in the config file or env
	if cfg.Database.DSN == "" {
		cfg.Database.DSN = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s connect_timeout=%d",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
			cfg.Database.Timezone,
			cfg.Database.ConnectTimeoutSeconds,
		)
	}

	return &cfg, nil
}
