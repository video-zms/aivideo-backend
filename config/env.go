package config

import (
	"github.com/sirupsen/logrus"
	"log"

	"github.com/BurntSushi/toml"
)

type (
	envConfig struct {
		Dev         bool
		LocalDev    bool
		LogConsole  bool
		ActLocalDev bool
		MainDBDsn   string
	}
)

var (
	EnvConfigFile = "./config/env.toml"
	envConf       envConfig
)

func InitEnvConf() {
	if _, err := toml.DecodeFile(EnvConfigFile, &envConf); err != nil {
		log.Fatal(err)
	}
}

func InitEnvConfCustom(file string) {
	if _, err := toml.DecodeFile(file, &envConf); err != nil {
		logrus.Fatal(err)
	}
}

func IsDev() bool {
	return envConf.Dev
}

func IsLocalDev() bool {
	return envConf.LocalDev
}

func IsActLocalDev() bool {
	return envConf.ActLocalDev
}

func GetMainMysqlDsn() string {
	return envConf.MainDBDsn
}
