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
	mut      sync.Mutex
}

func (conf *Config) Init() (err error) {
	if conf.omniplug != nil {
		return
	}

	conf.mut.Lock()
	defer conf.mut.Unlock()
	if conf.omniplug == nil {
		v, err := json.Marshal(conf)
		if err != nil {
			return err
		}

		conf.omniplug = &Omniplug{}
		if err = json.Unmarshal(v, conf.omniplug); err != nil {
			return err
		}

		if err = conf.omniplug.Init(); err != nil {
			return err
		}
	}
	return
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
