package config

import (
	"github.com/spf13/viper"
)

// Config結構
type Config struct {
	Group    string `mapstructure:"Group"`
	Port     string `mapstructure:"PORT"`
	DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri string `mapstructure:"REDIS_URL"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}
