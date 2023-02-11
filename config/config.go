package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App    AppConfig
		Server ServerConfig
		Logger LoggerConfig
		DB     string `yaml:"db"`
		Redis  *Redis
	}

	AppConfig struct {
		Mode       string `yaml:"mode"`
		LoginLimit int    `yaml:"loginLimit"`
		PassLimit  int    `yaml:"passLimit"`
		IpLimit    int    `yaml:"ipLimit"`
	}

	ServerConfig struct {
		Port              string        `yaml:"port"`
		PprofPort         string        `yaml:"pprofPort"`
		Mode              string        `yaml:"mode"`
		JwtSecretKey      string        `yaml:"jwtSecretKey"`
		CookieName        string        `yaml:"cookieName"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		SSL               bool          `yaml:"ssl"`
		CtxDefaultTimeout time.Duration `yaml:"ctxDefaultTimeout"`
		CSRF              bool          `yaml:"csrf"`
		Debug             bool          `yaml:"debug"`
		MaxConnectionIdle time.Duration `yaml:"maxConnectionIdle"`
		Timeout           time.Duration `yaml:"timeout"`
		MaxConnectionAge  time.Duration `yaml:"maxConnectionAge"`
		Time              time.Duration `yaml:"time"`
	}

	LoggerConfig struct {
		Development bool   `yaml:"development"`
		Level       string `yaml:"level"`
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}
)

// LoadConfig Load config file from given path -.
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
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
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	c.Redis = &Redis{
		Host:     v.GetString(envConfigRedisHost),
		Port:     v.GetString(envConfigRedisPort),
		Password: v.GetString(envConfigRedisPassword),
	}

	c.Redis.Host = v.GetString(envConfigRedisHost)
	if c.Redis.Host == "" {
		c.Redis.Host = "localhost"
	}
	c.Redis.Port = v.GetString(envConfigRedisPort)
	if c.Redis.Port == "" {
		c.Redis.Port = "6379"
	}
	c.Redis.Password = v.GetString(envConfigRedisPassword)

	return &c, nil
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
