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
		Mode     string `yaml:"Mode"`
		MaxLogin int    `yaml:"MaxLogin"`
		MaxPass  int    `yaml:"MaxPass"`
		MaxIp    int    `yaml:"MaxIP"`
	}

	ServerConfig struct {
		Port              string        `yaml:"Port"`
		PprofPort         string        `yaml:"PprofPort"`
		Mode              string        `yaml:"Mode"`
		JwtSecretKey      string        `yaml:"JwtSecretKey"`
		CookieName        string        `yaml:"CookieName"`
		ReadTimeout       time.Duration `yaml:"ReadTimeout"`
		WriteTimeout      time.Duration `yaml:"WriteTimeout"`
		SSL               bool          `yaml:"SSL"`
		CtxDefaultTimeout time.Duration `yaml:"CtxDefaultTimeout"`
		CSRF              bool          `yaml:"CSRF"`
		Debug             bool          `yaml:"Debug"`
		MaxConnectionIdle time.Duration `yaml:"MaxConnectionIdle"`
		Timeout           time.Duration `yaml:"Timeout"`
		MaxConnectionAge  time.Duration `yaml:"MaxConnectionAge"`
		Time              time.Duration `yaml:"Time"`
	}

	LoggerConfig struct {
		Development bool   `yaml:"Development"`
		Level       string `yaml:"Level"`
	}
)

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
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
