package scheduler

import (
	"fmt"
	"strings"

	"github.com/keenywheels/go-spy/internal/pkg/scraper"
	"github.com/keenywheels/go-spy/internal/scheduler/service"
	"github.com/spf13/viper"
)

// Config struct for logger config
type LoggerConfig struct {
	LogLevel      string `mapstructure:"loglvl"`
	Mode          string `mapstructure:"mode"`
	Encoding      string `mapstructure:"encoding"`
	LogPath       string `mapstructure:"log_path"`
	MaxLogSize    int    `mapstructure:"max_log_size"`
	MaxLogBackups int    `mapstructure:"max_log_backups"`
	MaxLogAge     int    `mapstructure:"max_log_age"`
}

// SystemServerConfig contains config for system server
type SystemServerConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

// AppConfig contains all configs which connected to main app
type AppConfig struct {
	CronPattern  string             `mapstructure:"cron_pattern"`
	WorkersCount int                `mapstructure:"workers_count"`
	Sites        []service.Site     `mapstructure:"sites"`
	LoggerCfg    LoggerConfig       `mapstructure:"logger"`
	ScraperCfg   scraper.Config     `mapstructure:"scraper"`
	SysSrvCfg    SystemServerConfig `mapstructure:"system_server"`
}

// KafkaTopics contains all kafka topics
type KafkaTopics struct {
	ScraperData string `mapstructure:"scraper_data"`
}

// KafkaConfig contains config for kafka
type KafkaConfig struct {
	MaxRetry int         `mapstructure:"max_retry"`
	Brokers  []string    `mapstructure:"brokers"`
	Topics   KafkaTopics `mapstructure:"topics"`
}

// Config global config, contains all configs
type Config struct {
	SchedulerCfg AppConfig   `mapstructure:"scheduler"`
	KafkaCfg     KafkaConfig `mapstructure:"kafka"`
}

// LoadConfig
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	// env support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
