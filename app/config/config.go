package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Viper *viper.Viper
	Env   Env
}

func NewConfig() *Config {

	viper.AutomaticEnv()
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	getViper := viper.GetViper()
	env := getEnv(getViper)

	return &Config{
		Viper: getViper,
		Env:   env,
	}
}

func getEnv(v *viper.Viper) Env {
	env := v.GetString("service-config.env")

	switch env {
	case string(AWS):
		return AWS
	case string(Local):
		return Local
	default:
		fmt.Printf("Error: config.env[%v] invalid. Using Default[local]", v.Get("config.env"))
		return Local
	}
}

type Env string

const (
	AWS   Env = "aws"
	Local Env = "local"
)
