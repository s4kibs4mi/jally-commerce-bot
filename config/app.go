package config

import (
	"github.com/spf13/viper"
)

// Application holds the application configuration
type Application struct {
	Name           string
	Host           string
	Port           int
	ShopemaaKey    string
	ShopemaaSecret string
	TwilioUsername string
	TwilioPassword string
	URL            string
}

// app is the default application configuration
var app Application

// App returns the default application configuration
func App() *Application {
	return &app
}

// LoadApp loads application configuration
func LoadApp() {
	mu.Lock()
	defer mu.Unlock()

	app = Application{
		Name:           viper.GetString("app.name"),
		Host:           viper.GetString("app.host"),
		Port:           viper.GetInt("app.port"),
		ShopemaaKey:    viper.GetString("app.shopemaa_key"),
		ShopemaaSecret: viper.GetString("app.shopemaa_secret"),
		TwilioUsername: viper.GetString("app.twilio_username"),
		TwilioPassword: viper.GetString("app.twilio_password"),
		URL:            viper.GetString("app.url"),
	}
}
