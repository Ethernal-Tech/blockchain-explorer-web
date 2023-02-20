package services

import (
	"fmt"
	"webbc/DB"
	"webbc/configuration"
	"webbc/eth"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/uptrace/bun"
)

type ConfigurationServiceImplementation struct {
	appConfig      *configuration.ApplicationConfiguration
	generalConfig  *configuration.GeneralConfiguration
	client         *rpc.Client
	database       *bun.DB
	addressService AddressService
}

func NewConfigurationService(appConfig *configuration.ApplicationConfiguration, config *configuration.GeneralConfiguration, client *rpc.Client, database *bun.DB, addressService AddressService) ConfigurationService {
	return &ConfigurationServiceImplementation{appConfig: appConfig, generalConfig: config, client: client, database: database, addressService: addressService}
}

func (csi *ConfigurationServiceImplementation) GetAppConfiguration() *configuration.ApplicationConfiguration {
	csi.appConfig.Mutex.RLock()
	defer csi.appConfig.Mutex.RUnlock()
	return csi.appConfig
}

func (csi *ConfigurationServiceImplementation) GetGeneralConfiguration() *configuration.GeneralConfiguration {
	csi.generalConfig.Mutex.RLock()
	defer csi.generalConfig.Mutex.RUnlock()
	return csi.generalConfig
}

func (csi *ConfigurationServiceImplementation) UpdateAppConfiguration(config *configuration.ApplicationConfiguration) error {
	csi.appConfig.Mutex.Lock()
	defer csi.appConfig.Mutex.Unlock()

	csi.appConfig.BackgroundColor = config.BackgroundColor
	csi.appConfig.MainColor = config.MainColor
	csi.appConfig.HeaderTitle = config.HeaderTitle
	csi.appConfig.MainTitle = config.MainTitle
	csi.appConfig.Viper.Set("BACKGROUND_COLOR", fmt.Sprintf("%q", csi.appConfig.BackgroundColor))
	csi.appConfig.Viper.Set("MAIN_COLOR", fmt.Sprintf("%q", csi.appConfig.MainColor))
	csi.appConfig.Viper.Set("HEADER_TITLE", fmt.Sprintf("%q", csi.appConfig.HeaderTitle))
	csi.appConfig.Viper.Set("MAIN_TITLE", fmt.Sprintf("%q", csi.appConfig.MainTitle))
	if err := csi.appConfig.Viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func (csi *ConfigurationServiceImplementation) UpdateGeneralConfiguration(config *configuration.GeneralConfiguration) error {
	csi.generalConfig.Mutex.Lock()
	defer csi.generalConfig.Mutex.Unlock()

	csi.client = eth.GetClient(config.HTTPUrl)
	csi.addressService.ChangeClient(csi.client)
	DB.ChangeDB(config, csi.database)

	csi.generalConfig.DataBaseHost = config.DataBaseHost
	csi.generalConfig.DataBasePort = config.DataBasePort
	csi.generalConfig.DataBaseName = config.DataBaseName
	csi.generalConfig.DataBaseUser = config.DataBaseUser
	csi.generalConfig.DataBasePassword = config.DataBasePassword
	csi.generalConfig.HTTPUrl = config.HTTPUrl

	csi.generalConfig.Viper.Set("DB_HOST", fmt.Sprintf("%q", csi.generalConfig.DataBaseHost))
	csi.generalConfig.Viper.Set("DB_PORT", fmt.Sprintf("%q", csi.generalConfig.DataBasePort))
	csi.generalConfig.Viper.Set("DB_NAME", fmt.Sprintf("%q", csi.generalConfig.DataBaseName))
	csi.generalConfig.Viper.Set("DB_USER", fmt.Sprintf("%q", csi.generalConfig.DataBaseUser))
	csi.generalConfig.Viper.Set("DB_PASSWORD", fmt.Sprintf("%q", csi.generalConfig.DataBasePassword))
	csi.generalConfig.Viper.Set("HTTPURL", fmt.Sprintf("%q", csi.generalConfig.HTTPUrl))
	if err := csi.generalConfig.Viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
