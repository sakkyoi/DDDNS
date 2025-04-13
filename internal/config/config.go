package config

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Mode string

const (
	ClientIpMode Mode = "ip"
	EcsMode      Mode = "ecs"
)

type Config struct {
	ListenHost   string `mapstructure:"listen_host"`
	DNSPort      int    `mapstructure:"dns_port"`
	ApiPort      int    `mapstructure:"api_port"`
	TTL          int    `mapstructure:"ttl"`
	Domain       string `mapstructure:"domain"`
	Mode         Mode   `mapstructure:"mode"`
	Fallback     string `mapstructure:"fallback"`
	FallbackType string `mapstructure:"fallback_type"`
	RedisHost    string `mapstructure:"redis_host"`
	RedisPort    int    `mapstructure:"redis_port"`
	RedisDB      int    `mapstructure:"redis_db"`
	RedisUser    string `mapstructure:"redis_user"`
	RedisPass    string `mapstructure:"redis_pass"`
	LogLevel     string `mapstructure:"log_level"`
}

func Load() *Config {
	viper.SetDefault("listen_host", "127.0.0.1")
	viper.SetDefault("dns_port", 53)
	viper.SetDefault("api_port", 8080)
	viper.SetDefault("ttl", 0)
	viper.SetDefault("domain", "")
	viper.SetDefault("mode", "ip")
	viper.SetDefault("fallback", "127.0.0.1")
	viper.SetDefault("fallback_type", "A")
	viper.SetDefault("redis_host", "")
	viper.SetDefault("redis_port", "")
	viper.SetDefault("redis_db", 0)
	viper.SetDefault("redis_user", "")
	viper.SetDefault("redis_pass", "")
	viper.SetDefault("log_level", "info")

	// Bind flags to viper
	pflag.String("listen_host", "127.0.0.1", "DDDNS server listen host")
	pflag.Int("dns_port", 53, "DNS server port")
	pflag.Int("api_port", 8080, "API server port")
	pflag.Int("ttl", 0, "TTL for DNS records in seconds, 0 for no expiration")
	pflag.String("domain", "", "Root domain for DDDNS")
	pflag.String("mode", "ip", "Mode for DDDNS (ip or subnet)")
	pflag.String("fallback", "127.0.0.1", "Fallback IP or CNAME")
	pflag.String("fallback_type", "A", "Fallback DNS record type (only A or CNAME is supported)")
	pflag.String("redis_host", "", "Redis server host")
	pflag.Int("redis_port", 6379, "Redis server port")
	pflag.Int("redis_db", 0, "Redis database number")
	pflag.String("redis_user", "", "Redis username")
	pflag.String("redis_pass", "", "Redis password")
	pflag.String("log_level", "info", "Log level (debug, info, warn, error, fatal)")
	pflag.String("config", "config.yaml", "Path to config file")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal("Error binding flags", "error", err)
	}

	// Environment variable binding
	viper.SetEnvPrefix("DDDNS")
	viper.AutomaticEnv()

	// Read in config file
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", "error", err)
	}

	// Unmarshal config into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Unable to decode into struct", "error", err)
	}

	if cfg.Domain == "" {
		log.Fatal("Domain is required")
	}

	if cfg.Mode != ClientIpMode && cfg.Mode != EcsMode {
		log.Fatal("Invalid mode, must be 'ip' or 'ecs'")
	}

	if cfg.FallbackType != "A" && cfg.FallbackType != "CNAME" {
		log.Fatal("Invalid fallback type, must be 'A' or 'CNAME'")
	}

	return &cfg
}
