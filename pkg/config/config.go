package config

import (
    "github.com/spf13/viper"
    "os"
)

type Config struct {
    ServerAddress string
    NasMountPath  string
    HostNasPath   string
    NasPath       string
    Port          string
}

func LoadConfig() (*Config, error) {
    if port := os.Getenv("PORT"); port != "" {
        viper.Set("PORT", port)
    }
    if nasPath := os.Getenv("NAS_PATH"); nasPath != "" {
        viper.Set("NAS_PATH", nasPath)
    }
    if nasMountPath := os.Getenv("NAS_MOUNT_PATH"); nasMountPath != "" {
        viper.Set("NAS_MOUNT_PATH", nasMountPath)
    }

    viper.AutomaticEnv()

    serverAddress := ":" + viper.GetString("PORT")

    return &Config{
        ServerAddress: serverAddress,
        NasMountPath:  viper.GetString("NAS_MOUNT_PATH"),
        HostNasPath:   viper.GetString("HOST_NAS_PATH"),
        NasPath:       viper.GetString("NAS_PATH"),
        Port:          viper.GetString("PORT"),
    }, nil
}
