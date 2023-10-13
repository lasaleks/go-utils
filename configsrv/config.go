package configsrv

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigSrv struct {
	TIME_ZONE string `yaml:"TimeZone"`

	WORKINGSHIFT struct {
		BEGINING_FIRST_SHIFT int `yaml:"BeginingFirstShift"`
		NUMBEROFSHIFT        int `yaml:"NumberOfShift"`
	} `yaml:"WorkingShift"`

	MYSQL struct {
		HOST     string `yaml:"host"`
		USER     string `yaml:"user"`
		PASSWORD string `yaml:"password"`
		PORT     int    `yaml:"port"`
		DATABASE string `yaml:"database"`
	} `yaml:"mysql"`

	REDIS struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`

	GRPC_SPPD struct {
		ADDRESS string `yaml:"address"`
	} `yaml:"grpc-sppd"`

	GRPC_SBEACON struct {
		ADDRESS string `yaml:"address"`
	} `yaml:"grpc-sbeacon"`

	GRPC_ALARMZONE struct {
		ADDRESS string `yaml:"address"`
	} `yaml:"grpc-alarmzone"`

	GRPC_SAUX struct {
		ADDRESS string `yaml:"address"`
	} `yaml:"grpc-saux"`

	GRPC_BLZONE struct {
		ADDRESS string `yaml:"address"`
	} `yaml:"grpc-blzone"`
}

func (conf *ConfigSrv) ParseConfig(pathfile string) error {
	yamlFile, err := os.ReadFile(pathfile)
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return fmt.Errorf("unmarshal: %s err   #%v ", pathfile, err)
	}
	tz := conf.TIME_ZONE
	if value, exists := os.LookupEnv("TIME_ZONE"); exists {
		tz = value
		log.Printf("TIME_ZONE:%s  from environment variables", tz)
	}
	if len(tz) == 0 {
		tz = "Asia/Novokuznetsk"
	}
	conf.TIME_ZONE = tz
	return nil
}

/*
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
*/
