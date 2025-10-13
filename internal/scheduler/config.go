package scheduler

import (
	"fmt"
	"strings"

	"github.com/keenywheels/go-spy/pkg/scraper"
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

// AppConfig contains all configs which connected to main app
type AppConfig struct {
	CronPattern string         `mapstructure:"cron_pattern"`
	Sites       []string       `mapstructure:"sites"`
	LoggerCfg   LoggerConfig   `mapstructure:"logger"`
	ScraperCfg  scraper.Config `mapstructure:"scraper"`
}

// Config global config, contains all configs
type Config struct {
	SchedulerCfg AppConfig `mapstructure:"scheduler"`
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
