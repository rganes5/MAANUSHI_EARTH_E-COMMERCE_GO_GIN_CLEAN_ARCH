package config

import (
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
}

var envs = []string{
	"DB_HOST", "DB_USER", "DB_NAME", "DB_PORT", "DB_PASSWORD", //DATABASE
	"JWT_CODE", //JWT
}

var config Config

func LoadConfig() (Config, error) {
	var config Config

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

	return config, nil
}

// to get the secret code for jwt
func GetJWTCofig() string {
	return config.JWT
}

func GetCofig() Config {
	return config
}
