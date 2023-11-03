package main

import (
	"errors"

	"github.com/Kong/go-pdk"
)

var (
	ErrDatasourceNoID         = errors.New("datasource's id can't be empty")
	ErrDatasourceIDNotAllowed = errors.New("datasource's id is invalid because the keyword is not allowed")
)

type Datasource struct {
	ID     string           `json:"id" yaml:"id"`
	Type   string           `json:"type" yaml:"type"`
	Static DatasourceStatic `json:"static" yaml:"static"`
	Kong   DatasourceKong   `json:"kong" yaml:"kong"`
	HTTP   DatasourceHTTP   `json:"http" yaml:"http"`
}

type configDatasource struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`

	Static configDatasourceStatic `json:"static" yaml:"static"`
	Kong   configDatasourceKong   `json:"kong" yaml:"kong"`
	HTTP   configDatasourceHTTP   `json:"http" yaml:"http"`
}

func (d *Datasource) Init() (err error) {
	if d.ID == "" {
		return ErrDatasourceNoID
	}
	if d.ID == "__self__" {
		return ErrDatasourceIDNotAllowed
	}

	switch d.Type {
	case "static":
		return d.Static.Init()
	case "kong":
		return d.Kong.Init()
	case "http":
		return d.HTTP.Init()
	}
	return
}

func (d *Datasource) ToData(kong *pdk.PDK, data Data) (m map[string]interface{}, err error) {
	switch d.Type {
	case "static":
		return d.Static.ToData(kong, data)
	case "kong":
		return d.Kong.ToData(kong, data)
	case "http":
		return d.HTTP.ToData(kong, data)
	}
	return
}
