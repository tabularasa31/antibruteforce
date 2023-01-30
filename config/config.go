package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App    AppConfig
		Server ServerConfig
		Logger LoggerConfig
	}

	AppConfig struct {
		Mode     string `yaml:"mode"`
		MaxLogin int    `yaml:"maxLogin"`
		MaxPass  int    `yaml:"maxPass"`
		MaxIP    int    `yaml:"maxIp"`
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
)

// Load config file from given path -.
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
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

// Parse config file -.
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config -.
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
