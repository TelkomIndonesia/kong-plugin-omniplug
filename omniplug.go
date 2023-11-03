package main

import (
	"fmt"

	"github.com/Kong/go-pdk"
)

type Data map[string]interface{}

type Omniplug struct {
	Datasources []*Datasource `json:"datasources" yaml:"datasources"`
	Injector    Injector      `json:"injector" yaml:"injector"`
}

func (r *Omniplug) Init() (err error) {
	m := map[string]struct{}{}
	for _, d := range r.Datasources {
		if err = d.Init(); err != nil {
			return
		}
		if _, ok := m[d.ID]; ok {
			return fmt.Errorf("duplicate datasource id: %s", d.ID)
		}
		m[d.ID] = struct{}{}
	}
	return
}

func (r *Omniplug) Access(kong *pdk.PDK) {
	data, err := r.ResolvData(kong)
	if err != nil {
		kong.Response.Exit(500, []byte(err.Error()), nil)
	}

	err = r.Injector.Inject(kong, data)
	if err != nil {
		kong.Response.Exit(500, []byte(err.Error()), nil)
	}
}

func (r *Omniplug) ResolvData(kong *pdk.PDK) (data Data, err error) {
	data = map[string]interface{}{}
	for _, d := range r.Datasources {
		data[d.ID], err = d.ToData(kong, data)
		if err != nil {
			return data, err
		}
	}
	return
}
