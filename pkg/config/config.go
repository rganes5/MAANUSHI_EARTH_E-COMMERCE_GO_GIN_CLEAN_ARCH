package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	JWT        string `mapstructure:"JWT_CODE"`
	AUTHTOKEN  string `mapstructure:"AUTH_TOKEN"`
	ACCOUNTSID string `mapstructure:"ACCOUNT_SID"`
	SERVICESID string `mapstructure:"SERVICE_SID"`
}

var envs = []string{
	"DB_HOST", "DB_USER", "DB_NAME", "DB_PORT", "DB_PASSWORD", //DATABASE
	"JWT_CODE",                                 //JWT
	"AUTH_TOKEN", "ACCOUNT_SID", "SERVICE_SID", //twilio details
}

var config Config

func LoadConfig() (Config, error) {

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}
	fmt.Println("config", config)
	return config, nil
}

// to get the secret code for jwt
func GetJWTCofig() string {
	return config.JWT
}

func GetCofig() Config {
	return config
}
