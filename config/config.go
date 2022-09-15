package config

import (
	"os"
	"sync"

	"github.com/spf13/viper"
	// this package is necessary to read config from remote consul
	_ "github.com/spf13/viper/remote"
)

var mu sync.Mutex

// LoadConfig initiates of config load
func LoadConfig() error {
	return LoadConfigWithPath(os.Getenv("CONFIG_DIR"))
}

func LoadConfigWithPath(path string) error {
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	LoadApp()

	return nil
}
