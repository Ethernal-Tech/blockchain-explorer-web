package configuration

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

type GeneralConfiguration struct {
	Mutex                         *sync.RWMutex
	DataBaseUser                  string `form:"dbUser"`
	DataBasePassword              string `form:"dbPassword"`
	DataBaseHost                  string `form:"dbHost"`
	DataBasePort                  string `form:"dbPort"`
	DataBaseName                  string `form:"dbName"`
	HTTPUrl                       string `form:"nodeUrl"`
	CallTimeoutInSeconds          uint
	TransactionsMaxCount          uint64
	TransactionsByAddressMaxCount uint
	NftLatestTransfersMaxCount    uint
	Viper                         *viper.Viper
}

type ApplicationConfiguration struct {
	Mutex           *sync.RWMutex
	BackgroundColor string `form:"backgroundColor"`
	MainColor       string `form:"mainColor"`
	HeaderTitle     string `form:"headerTitle"`
	MainTitle       string `form:"mainTitle"`
	Viper           *viper.Viper
}

type AuthConfiguration struct {
	Mutex    *sync.RWMutex
	Username string
	Password string
	Viper    *viper.Viper
}

func LoadConfiguration() (*GeneralConfiguration, *ApplicationConfiguration, *AuthConfiguration) {
	generalConfiguration := GeneralConfiguration{
		Viper: getNewViperInstance("general.env"),
		Mutex: &sync.RWMutex{},
	}
	generalConfiguration.DataBaseHost = generalConfiguration.Viper.GetString("DB_HOST")
	generalConfiguration.DataBasePort = generalConfiguration.Viper.GetString("DB_PORT")
	generalConfiguration.DataBaseName = generalConfiguration.Viper.GetString("DB_NAME")
	generalConfiguration.DataBaseUser = generalConfiguration.Viper.GetString("DB_USER")
	generalConfiguration.DataBasePassword = generalConfiguration.Viper.GetString("DB_PASSWORD")
	generalConfiguration.HTTPUrl = generalConfiguration.Viper.GetString("HTTPUrl")
	generalConfiguration.CallTimeoutInSeconds = generalConfiguration.Viper.GetUint("CALL_TIMEOUT_IN_SECONDS")
	generalConfiguration.TransactionsMaxCount = generalConfiguration.Viper.GetUint64("TRANSACTIONS_MAX_COUNT")
	generalConfiguration.TransactionsByAddressMaxCount = generalConfiguration.Viper.GetUint("TRANSACTIONS_BY_ADDRESS_MAX_COUNT")
	generalConfiguration.NftLatestTransfersMaxCount = generalConfiguration.Viper.GetUint("NFT_LATEST_TRANSFER_MAX_COUNT")

	appConfiguration := ApplicationConfiguration{
		Viper: getNewViperInstance("app.env"),
		Mutex: &sync.RWMutex{},
	}
	appConfiguration.BackgroundColor = appConfiguration.Viper.GetString("BACKGROUND_COLOR")
	appConfiguration.MainColor = appConfiguration.Viper.GetString("MAIN_COLOR")
	appConfiguration.HeaderTitle = appConfiguration.Viper.GetString("HEADER_TITLE")
	appConfiguration.MainTitle = appConfiguration.Viper.GetString("MAIN_TITLE")

	authConfiguration := AuthConfiguration{
		Viper: getNewViperInstance("auth.env"),
		Mutex: &sync.RWMutex{},
	}
	authConfiguration.Username = authConfiguration.Viper.GetString("USERNAME")
	authConfiguration.Password = authConfiguration.Viper.GetString("PASSWORD")

	return &generalConfiguration, &appConfiguration, &authConfiguration
}

func getNewViperInstance(fileName string) *viper.Viper {
	configFile, error := filepath.Abs(fileName)

	if error != nil {
		//TODO: panic
		os.Exit(-1)
	}
	newViper := viper.New()
	newViper.SetConfigFile(configFile)
	newViper.SetConfigType("env")
	error = newViper.ReadInConfig()

	if error != nil {
		//TODO: panic
		os.Exit(-1)
	}
	return newViper
}
