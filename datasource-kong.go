package main

import (
	"strings"

	"github.com/Kong/go-pdk"
)

type DatasourceKong struct {
	Consumer           bool      `json:"consumer" yaml:"consumer"`
	Context            []string  `json:"context" yaml:"context"`
	RequestHeaders     *[]string `json:"request_headers" yaml:"request_headers"`
	RequestQueryParams *[]string `json:"request_query_params" yaml:"request_query_params"`
}

type configDatasourceKong DatasourceKong

func (dk DatasourceKong) Init() (err error) {
	return
}

func (dk DatasourceKong) ToData(kong *pdk.PDK, data Data) (m map[string]interface{}, err error) {
	m = map[string]interface{}{}

	m["consumer"], err = dk.getConsumer(kong)
	if err != nil {
		return
	}

	m["context"], err = dk.getContext(kong)
	if err != nil {
		return
	}

	m["request_headers"], err = dk.getHeaders(kong)
	if err != nil {
		return
	}

	m["request_query_params"], err = dk.getQueryParams(kong)
	return
}

func (dk DatasourceKong) getConsumer(kong *pdk.PDK) (user map[string]interface{}, err error) {
	if !dk.Consumer {
		return
	}

	c, err := kong.Client.GetConsumer()
	if err != nil {
		return
	}
	user = map[string]interface{}{
		"id":         c.Id,
		"custom_id":  c.CustomId,
		"username":   c.Username,
		"created_at": c.CreatedAt,
		"tags":       c.Tags,
	}
	return
}
func (dk DatasourceKong) getContext(kong *pdk.PDK) (ctx map[string]interface{}, err error) {
	if len(dk.Context) == 0 {
		return
	}

	ctx = map[string]interface{}{}
	for _, k := range dk.Context {
		ctx[k], err = kong.Ctx.GetSharedAny(k)
		if err != nil {
			return
		}
	}
	return
}

func (dk DatasourceKong) getHeaders(kong *pdk.PDK) (h map[string][]string, err error) {
	if dk.RequestHeaders == nil {
		return
	}

	if len(*dk.RequestHeaders) == 0 {
		h, err = kong.Request.GetHeaders(-1)
		if err != nil {
			return
		}
		for k, v := range h {
			delete(h, k)
			h[strings.ReplaceAll(strings.ToLower(k), "-", "_")] = v
		}
		return
	}

	h = map[string][]string{}
	for _, k := range *dk.RequestHeaders {
		v, err := kong.Request.GetHeader(k)
		if err != nil {
			return nil, err
		}
		h[strings.ReplaceAll(strings.ToLower(k), "-", "_")] = []string{v}
	}
	return
}

func (dk DatasourceKong) getQueryParams(kong *pdk.PDK) (q map[string][]string, err error) {
	if dk.RequestQueryParams == nil {
		return
	}

	if len(*dk.RequestQueryParams) == 0 {
		q, err = kong.Request.GetQuery(-1)
		if err != nil {
			return
		}
		return
	}

	q = map[string][]string{}
	for _, k := range *dk.RequestQueryParams {
		v, err := kong.Request.GetQueryArg(k)
		if err != nil {
			return nil, err
		}
		q[k] = []string{v}
	}
	return
}
