package util

import (
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    DBDriver string `mapstructure:"DB_DRIVER"`
    DBSource string `mapstructure:"DB_SOURCE"`
    ServerAddres string `mapstructure:"SERVICE_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env") // we can use json, html or any one if we want

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    log.Println("Config file:", viper.ConfigFileUsed())
    if err != nil {
        log.Printf("ReadInConfig error: %v", err)
        return
    }

    viper.Unmarshal(&config)
    return
}