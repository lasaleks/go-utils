package configsrv

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		return fmt.Errorf("configsrv yamlFile.Get err   #%v, pathfile:%s", err, pathfile)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return fmt.Errorf("configsrv unmarshal: %s err   #%v ", pathfile, err)
	}
	tz := conf.TIME_ZONE
	if value, exists := os.LookupEnv("TIME_ZONE"); exists {
		tz = value
		log.Printf("TIME_ZONE:%s read from environment variables TIME_ZONE", tz)
	}
	if len(tz) == 0 {
		tz = "Asia/Novokuznetsk"
	}
	conf.TIME_ZONE = tz
	if value, exists := os.LookupEnv("WORKING_SHIFT_BEGIN_FIRST"); exists {
		if v, err := strconv.ParseInt(value, 10, 32); err == nil {
			conf.WORKINGSHIFT.BEGINING_FIRST_SHIFT = int(v)
			log.Printf("WorkingShift.BeginingFirstShift:%d read from environment variables WORKING_SHIFT_BEGIN_FIRST", v)
		} else {
			log.Printf("Benvironment variables WORKING_SHIFT_BEGIN_FIRST error %s", err)
		}
	}
	if value, exists := os.LookupEnv("WORKING_SHIFT_NOF"); exists {
		if v, err := strconv.ParseInt(value, 10, 32); err == nil {
			conf.WORKINGSHIFT.NUMBEROFSHIFT = int(v)
			log.Printf("WorkingShift.NumberOfShift:%d read from environment variables WORKING_SHIFT_NOF", v)
		} else {
			log.Printf("Benvironment variables WORKING_SHIFT_NOF error %s", err)
		}
	}
	/*
		TIME_ZONE
		WORKING_SHIFT_BEGIN_FIRST
		WORKING_SHIFT_NOF
	*/
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
