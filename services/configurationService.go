package services

import (
	"fmt"
	"webbc/configuration"

	"github.com/spf13/viper"
)

type ConfigurationServiceImplementation struct {
	config *configuration.TemplateConfiguration
}

func NewConfigurationService(templateConfig *configuration.TemplateConfiguration) ConfigurationService {
	return &ConfigurationServiceImplementation{config: templateConfig}
}

func (csi *ConfigurationServiceImplementation) Update(config *configuration.TemplateConfiguration) error {
	csi.config.Mutex.Lock()
	defer csi.config.Mutex.Unlock()

	csi.config.BackgroundColor = config.BackgroundColor
	csi.config.MainColor = config.MainColor
	csi.config.HeaderTitle = config.HeaderTitle
	csi.config.MainTitle = config.MainTitle

	viper.Set("BACKGROUND_COLOR", fmt.Sprintf("%q", csi.config.BackgroundColor))
	viper.Set("MAIN_COLOR", fmt.Sprintf("%q", csi.config.MainColor))
	viper.Set("HEADER_TITLE", fmt.Sprintf("%q", csi.config.HeaderTitle))
	viper.Set("MAIN_TITLE", fmt.Sprintf("%q", csi.config.MainTitle))
	if err := viper.WriteConfig(); err != nil {
	}
	return nil
}
