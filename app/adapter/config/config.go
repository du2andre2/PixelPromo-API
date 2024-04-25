package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Viper *viper.Viper
}

func NewConfig() *Config {

	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Erro fatal no arquivo de configuração: %w \n", err))
	}

	return &Config{
		Viper: viper.GetViper(),
	}
}

type File struct{}
