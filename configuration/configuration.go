package configuration

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Configuration struct {
	DataBaseUser     string
	DataBasePassword string
	DataBaseHost     string
	DataBasePort     string
	DataBaseName     string
}

func LoadConfiguration() *Configuration {
	configurationFile, error := filepath.Abs(".env")

	if error != nil {
		//TODO: panic
		os.Exit(-1)
	}

	viper.SetConfigFile(configurationFile)
	error = viper.ReadInConfig()

	if error != nil {
		//TODO: panic
		os.Exit(-1)
	}

	configuration := Configuration{
		DataBaseUser:     viper.GetString("DB_USER"),
		DataBasePassword: viper.GetString("DB_PASSWORD"),
		DataBaseHost:     viper.GetString("DB_HOST"),
		DataBasePort:     viper.GetString("DB_PORT"),
		DataBaseName:     viper.GetString("DB_NAME"),
	}

	return &configuration
}
