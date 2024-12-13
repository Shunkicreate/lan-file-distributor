package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    ServerAddress string
    NasMountPath  string
}

func LoadConfig() (*Config, error) {
    viper.SetDefault("SERVER_ADDRESS", ":8080")
    viper.SetDefault("NAS_MOUNT_PATH", "/mnt/nas")
    viper.AutomaticEnv()

    return &Config{
        ServerAddress: viper.GetString("SERVER_ADDRESS"),
        NasMountPath:  viper.GetString("NAS_MOUNT_PATH"),
    }, nil
}
