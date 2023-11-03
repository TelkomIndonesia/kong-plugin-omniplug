package main

import (
	"github.com/Kong/go-pdk"
	"gopkg.in/yaml.v3"
)

type DatasourceStatic struct {
	i    map[string]interface{}
	text []byte
}

type configDatasourceStatic string

func (ds DatasourceStatic) Init() (err error) {
	return
}

func (ds DatasourceStatic) ToData(kong *pdk.PDK, data Data) (k map[string]interface{}, err error) {
	return ds.i, nil
}

func (ds *DatasourceStatic) UnmarshalText(text []byte) (err error) {
	ds.i = map[string]interface{}{}
	if err = yaml.Unmarshal(text, &ds.i); err != nil {
		return err
	}
	ds.text = text
	return
}

func (ds DatasourceStatic) MarshalText() (text []byte, err error) {
	return ds.text, nil
}
