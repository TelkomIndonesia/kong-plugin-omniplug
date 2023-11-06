package main

import (
	"encoding/json"
	"sync"

	"github.com/Kong/go-pdk"
)

type Config struct {
	Datasources []configDatasource `json:"datasources" yaml:"datasources"`
	Injector    configInjector     `json:"injector" yaml:"injector"`

	omniplug *Omniplug
	once     sync.Once
	onceErr  error
}

func (conf *Config) Init() (err error) {
	conf.once.Do(func() {
		v, err := json.Marshal(conf)
		if err != nil {
			conf.onceErr = err
			return
		}

		conf.omniplug = &Omniplug{}
		if conf.onceErr = json.Unmarshal(v, conf.omniplug); conf.onceErr != nil {
			return
		}

		if conf.onceErr = conf.omniplug.Init(); conf.onceErr != nil {
			return
		}
	})
	return conf.onceErr
}

func (conf *Config) Access(kong *pdk.PDK) {
	if err := conf.Init(); err != nil {
		kong.Log.Err(err.Error())
		kong.Response.ExitStatus(500)
	}

	conf.omniplug.Access(kong)
}

func NewConfig() interface{} {
	return &Config{}
}
