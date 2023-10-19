package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Kong/go-pdk"
	"github.com/Masterminds/sprig/v3"
)

type Data []map[string]interface{}

type ReqEnrich struct {
	Datasource []Datasource `json:"datasource"`
	Injector   Injector     `json:"injector"`
}

func New() interface{} {
	return &ReqEnrich{}
}

func (conf *ReqEnrich) Access(kong *pdk.PDK) {
	data, err := conf.ResolvData(kong)
	if err != nil {
		kong.Response.Exit(500, []byte(err.Error()), nil)
	}

	err = conf.Injector.Inject(kong, data)
	if err != nil {
		kong.Response.Exit(500, []byte(err.Error()), nil)
	}
}

func (conf *ReqEnrich) ResolvData(kong *pdk.PDK) (data Data, err error) {
	for _, d := range conf.Datasource {
		switch d.Type {
		case "static":
			data = append(data, d.Static)
		case "kong":
			k, err := d.Kong.ToData(kong)
			if err != nil {
				return nil, err
			}
			data = append(data, k)
		}
	}
	return
}

type Datasource struct {
	Type   string                 `json:"type"`
	Static map[string]interface{} `json:"static"`
	Kong   DatasourceKong         `json:"kong"`
}

type DatasourceKong struct {
	Consumer           bool `json:"consumer"`
	RequestHeaders     bool `json:"request_headers"`
	RequestQueryParams bool `json:"request_query_params"`
}

func (dk DatasourceKong) ToData(kong *pdk.PDK) (k map[string]interface{}, err error) {
	k = map[string]interface{}{}
	if dk.Consumer {
		k["consumer"], err = kong.Client.GetConsumer()
		if err != nil {
			return
		}
	}
	if dk.RequestHeaders {
		m, err := kong.Request.GetHeaders(-1)
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			delete(m, k)
			m[strings.ReplaceAll(strings.ToLower(k), "-", "_")] = v
		}
		k["request_headers"] = m
	}
	if dk.RequestQueryParams {
		k["request_query_params"], err = kong.Request.GetQuery(-1)
		if err != nil {
			return
		}
	}
	return
}

type Injector struct {
	Headers     map[string]Template `json:"headers"`
	QueryParams map[string]Template `json:"query_params"`
}

func (i Injector) Inject(kong *pdk.PDK, data Data) (err error) {
	d := map[string]interface{}{"data": data}

	h := map[string][]string{}
	for k, v := range i.Headers {
		s, err := v.ToString(d)
		if err != nil {
			return err
		}
		h[k] = []string{s}
	}
	if len(h) > 0 {
		kong.ServiceRequest.SetHeaders(h)
	}

	q := map[string][]string{}
	for k, v := range i.QueryParams {
		s, err := v.ToString(d)
		if err != nil {
			return err
		}
		h[k] = []string{s}
	}
	if len(q) > 0 {
		kong.ServiceRequest.SetQuery(q)
	}

	return
}

type Template struct {
	*template.Template

	text []byte
}

func NewTemplate(s string) (Template, error) {
	t := Template{}
	return t, t.UnmarshalText([]byte(s))
}

func (t *Template) UnmarshalText(text []byte) (err error) {
	t.text = text
	t.Template, err = template.New("template").Funcs(sprig.FuncMap()).Parse(string(text))
	return
}

func (t Template) MarshalText() (text []byte, err error) {
	return t.text, nil
}

func (t Template) ToString(data interface{}) (s string, err error) {
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return
	}
	return buf.String(), nil
}
