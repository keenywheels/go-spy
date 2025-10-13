package webapp

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// HttpConfig config for http server
type HttpConfig struct {
	Port            string        `mapstructure:"port"`
	Host            string        `mapstructure:"host"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

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

// S2SClient contains s2s client info
type S2SClient struct {
	Name  string `mapstructure:"name"`
	Token string `mapstructure:"token"`
}

// S2SConfig contains s2s info
type S2SConfig struct {
	Header  string      `mapstructure:"header"`
	Clients []S2SClient `mapstructure:"clients"`
}

// AppConfig contains all configs which connected to main app
type AppConfig struct {
	HttpCfg   HttpConfig   `mapstructure:"http"`
	LoggerCfg LoggerConfig `mapstructure:"logger"`
	S2SCfg    S2SConfig    `mapstructure:"s2s"`
}

// Config global config, contains all configs
type Config struct {
	AppCfg AppConfig `mapstructure:"app"`
}

// LoadConfig function which reads config file and return Config instance
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
