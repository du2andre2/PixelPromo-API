package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Viper *viper.Viper
	Env   Env
}

func NewConfig() *Config {

	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Erro fatal no arquivo de configuração: %w \n", err))
	}
	viper := viper.GetViper()
	env := getEnv(viper)

	return &Config{
		Viper: viper,
		Env:   env,
	}
}

func getEnv(v *viper.Viper) Env {
	env := v.GetString("service-config.env")

	switch env {
	case string(Dev):
		return Dev
	case string(Local):
		return Local
	default:
		fmt.Printf("Error: config.env[%v] invalid. Using Default[local]", v.Get("config.env"))
		return Local
	}
}

type Env string

const (
	Dev   Env = "dev"
	Local Env = "local"
)
