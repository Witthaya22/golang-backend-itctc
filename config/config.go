package config

import (
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   *Server   `mapstructure:"server" validate:"required"`
		Database *Database `mapstructure:"database" validate:"required"`
		Jwt      *Jwt      `mapstructure:"jwt" validate:"required"`
	}

	Server struct {
		Host         string        `mapstructure:"host" validate:"required"`
		Port         int           `mapstructure:"port" validate:"required"`
		Name         string        `mapstructure:"name" validate:"required"`
		AllowOrigins []string      `mapstructure:"allowOrigins" validate:"required"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout" validate:"required"`
		Timeout      time.Duration `mapstructure:"timeout" validate:"required"`
		BodyLimit    int64         `mapstructure:"bodyLimit" validate:"required"`
		Version      string        `mapstructure:"version" validate:"required"`
	}

	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		Dbname   string `mapstructure:"dbname" validate:"required"`
		Sslmode  string `mapstructure:"sslmode" validate:"required"`
		Schema   string `mapstructure:"schema" validate:"required"`
	}

	Jwt struct {
		SecretKey      string        `mapstructure:"secretKey" validate:"required"`
		adminKey       string        `mapstructure:"adminKey" validate:"required"`
		apiKey         string        `mapstructure:"apiKey" validate:"required"`
		AccessExpires  time.Duration `mapstructure:"accessExpires" validate:"required"`
		RefreshExpires time.Duration `mapstructure:"refreshExpires" validate:"required"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func ConfigGeting() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}

		validating := validator.New()
		if err := validating.Struct(configInstance); err != nil {
			panic(err)
		}
	})

	return configInstance
}
