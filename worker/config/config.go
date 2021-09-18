package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

var Config config

type config struct {
	Logger struct {
		Project     string   `yaml:"project"`
		App         string   `yaml:"app"`
		File        string   `yaml:"file"`
		Level       string   `yaml:"level"`
		Outputs     []string `yaml:"outputs,flow"`
		Hooks       []string `yaml:"hooks,flow"`
		EsServer    []string `yaml:"es_server,flow"`
		StashServer string   `yaml:"stash_server"`
	} `yaml:"logger"`

	Etcd struct {
		Endpoints   []string      `yaml:"endpoints,flow"`
		DialTimeout time.Duration `yaml:"dial_timeout"`
	} `yaml:"etcd"`
}

func InitConfig(filepath string) error {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, &Config)
}
