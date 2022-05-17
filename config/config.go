package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"sync"
	"time"
)

type Config struct {
	Portal struct {
		Listen string
	}
	Db struct {
		Host    string
		Port    int
		User    string
		Pass    string
		Schema  string
		Charset string
	}
	Email struct {
		User string
		Pass string
		Host string
		Port int
		Sign string
	}
	EmailValidation struct {
		ExpireDuration string `yaml:"expire_duration"`
	} `yaml:"email_validation"`
	Auth struct {
		Session struct {
			ExpireSeconds       time.Duration `yaml:"expire_seconds"`
			MaximumRenewalTimes int           `yaml:"maximum_renewal_times"`
		}
		Jwt struct {
			Issuer        string
			Secret        string
			ExpireSeconds time.Duration `yaml:"expire_seconds"`
		}
		Signature struct {
			TimestampToleranceSeconds time.Duration `yaml:"timestamp_tolerance_seconds"`
			BodyDigestTruncateBytes   int           `yaml:"body_digest_truncate_bytes"`
		}
	}
}

var conf Config
var confOnce sync.Once

func GetConfig() Config {
	confOnce.Do(func() {
		data := readConfigFile("configs.yaml")
		err := yaml.Unmarshal(data, &conf)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	})
	return conf
}
