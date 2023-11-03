package main

import "github.com/Kong/go-pdk"

type Injector struct {
	Context     map[string]map[string]Template `json:"context" yaml:"contexts"`
	Headers     map[string]Template            `json:"headers" yaml:"headers"`
	QueryParams map[string]Template            `json:"query_params" yaml:"query_params"`
}

type configInjector struct {
	Context     map[string]map[string]string `json:"context" yaml:"context"`
	Headers     map[string]string            `json:"headers" yaml:"headers"`
	QueryParams map[string]string            `json:"query_params" yaml:"query_params"`
}

func (i Injector) Inject(kong *pdk.PDK, data Data) (err error) {
	for k, v := range i.Context {
		m := map[string]interface{}{}
		for k1, v1 := range v {
			s, err := v1.ToString(data)
			if err != nil {
				return err
			}
			m[k1] = s
		}
		kong.Ctx.SetShared(k, m)
	}

	h := map[string][]string{}
	for k, v := range i.Headers {
		s, err := v.ToString(data)
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
		s, err := v.ToString(data)
		if err != nil {
			return err
		}
		q[k] = []string{s}
	}
	if len(q) > 0 {
		kong.ServiceRequest.SetQuery(q)
	}

	return
}
