package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      AppConfig
		Server   ServerConfig
		Logger   LoggerConfig
		Redis    Redis
		Postgres Postgres
	}

	AppConfig struct {
		Mode       string `yaml:"mode"`
		LoginLimit int    `yaml:"loginLimit"`
		PassLimit  int    `yaml:"passLimit"`
		IPLimit    int    `yaml:"ipLimit"`
	}

	ServerConfig struct {
		Port              string        `yaml:"port"`
		Mode              string        `yaml:"mode"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		SSL               bool          `yaml:"ssl"`
		CtxDefaultTimeout time.Duration `yaml:"ctxDefaultTimeout"`
		Debug             bool          `yaml:"debug"`
	}

	LoggerConfig struct {
		Development bool   `yaml:"development"`
		Level       string `yaml:"level"`
	}

	Redis struct {
		Host     string `env:"APP_REDIS_HOST"`
		Port     string `env:"APP_REDIS_PORT"`
		Password string `env:"APP_REDIS_PASSWORD"`
	}

	Postgres struct {
		Dsn     string `yaml:"dsn"`
		PoolMax int    `yaml:"poolMax" env:"PG_POOL_MAX"`
	}
)

// LoadConfig Load config file from given path -.
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(filename)
	v.AddConfigPath(".")
	v.SetEnvPrefix("app")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		var notFoundError *viper.ConfigFileNotFoundError
		if errors.As(err, &notFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// ParseConfig Parse config file -.
func ParseConfig(v *viper.Viper) (*Config, error) {
	c := &Config{Redis: Redis{"", "", ""}}

	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	c.Redis.Host = v.GetString(envConfigRedisHost)
	if c.Redis.Host == "" {
		v.SetDefault(envConfigRedisHost, "localhost")
		c.Redis.Host = v.GetString(envConfigRedisHost)
	}
	c.Redis.Port = v.GetString(envConfigRedisPort)
	if c.Redis.Port == "" {
		c.Redis.Port = "6379"
	}
	c.Redis.Password = v.GetString(envConfigRedisPassword)

	return c, nil
}

// GetConfig Get config -.
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
