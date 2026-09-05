package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Enviroment           string        `mapstructure:"ENVIROMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddres     string        `mapstructure:"HTTP_SERVICE_ADDRESS"`
	GRPCServerAddres     string        `mapstructure:"GRPC_SERVICE_ADDRESS"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // we can use json, html or any one if we want

	viper.AutomaticEnv()

	// Bind every key explicitly so values can come purely from environment
	// variables (e.g. in CI) even when no app.env file is present.
	for _, key := range []string{
		"DB_DRIVER",
		"DB_SOURCE",
		"SERVICE_ADDRESS",
		"REDIS_ADDRESS",
		"TOKEN_SYMMETRIC_KEY",
		"ACCESS_TOKEN_DURATION",
		"REFRESH_TOKEN_DURATION",
	} {
		if bindErr := viper.BindEnv(key); bindErr != nil {
			return config, bindErr
		}
	}

	err = viper.ReadInConfig()
	log.Println("Config file:", viper.ConfigFileUsed())
	if err != nil {
		// A missing config file is fine as long as the environment supplies
		// the required values; only a malformed file is a real error.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("ReadInConfig error: %v", err)
			return
		}
		err = nil
	}

	err = viper.Unmarshal(&config)
	return
}
