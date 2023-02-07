package configuration

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

type Configuration struct {
	DataBaseUser         string
	DataBasePassword     string
	DataBaseHost         string
	DataBasePort         string
	DataBaseName         string
	HTTPUrl              string
	CallTimeoutInSeconds uint
}

type TemplateConfiguration struct {
	Mutex           *sync.RWMutex
	BackgroundColor string `form:"backgroundColor"`
	MainColor       string `form:"mainColor"`
	HeaderTitle     string `form:"headerTitle"`
	MainTitle       string `form:"mainTitle"`
}

func LoadConfiguration() (*Configuration, *TemplateConfiguration) {
	configurationFile, error := filepath.Abs(".env")

	if error != nil {
		//TODO: panic
		os.Exit(-1)
	}

	viper.SetConfigFile(configurationFile)
	viper.SetConfigType("env")
	error = viper.ReadInConfig()

	if error != nil {
		//TODO: panic
		os.Exit(-1)
	}

	configuration := Configuration{
		DataBaseUser:         viper.GetString("DB_USER"),
		DataBasePassword:     viper.GetString("DB_PASSWORD"),
		DataBaseHost:         viper.GetString("DB_HOST"),
		DataBasePort:         viper.GetString("DB_PORT"),
		DataBaseName:         viper.GetString("DB_NAME"),
		HTTPUrl:              viper.GetString("HTTPUrl"),
		CallTimeoutInSeconds: viper.GetUint("CALL_TIMEOUT_IN_SECONDS"),
	}

	templateConfiguration := TemplateConfiguration{
		Mutex:           &sync.RWMutex{},
		BackgroundColor: viper.GetString("BACKGROUND_COLOR"),
		MainColor:       viper.GetString("MAIN_COLOR"),
		HeaderTitle:     viper.GetString("HEADER_TITLE"),
		MainTitle:       viper.GetString("MAIN_TITLE"),
	}

	return &configuration, &templateConfiguration
}
