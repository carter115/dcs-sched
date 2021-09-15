package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

var Config config

type config struct {
	Server struct {
		Addr           string        `yaml:"addr"`
		ReadTimeout    time.Duration `yaml:"read_timeout"`
		WriteTimeout   time.Duration `yaml:"write_timeout"`
		MaxHeaderBytes int           `yaml:"max_header_bytes"`
	} `yaml:"server"`

	Logger struct {
		File  []string `yaml:"file,flow"`
		Level string   `yaml:"level"`
	} `yaml:"logger"`

	Etcd struct {
		Endpoints   []string      `yaml:"endpoints,flow"`
		DialTimeout time.Duration `yaml:"dial_timeout"`
	} `yaml:"etcd"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
}

func InitConfig(filepath string) error {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, &Config)
}
