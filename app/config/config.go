package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Viper *viper.Viper
	Env   Env
}

func NewConfig() *Config {

	setupViper("./env.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat("./cm.yml"); err == nil {
		setupViper("./cm.yml")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
	}

	getViper := viper.GetViper()
	env := getEnv(getViper)

	return &Config{
		Viper: getViper,
		Env:   env,
	}
}

func setupViper(file string) {
	viper.SetConfigFile(file)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
