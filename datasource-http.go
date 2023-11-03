package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/dgraph-io/ristretto"
	"github.com/spf13/cast"
)

type Duration time.Duration

func (t *Duration) UnmarshalText(text []byte) (err error) {
	d, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*t = Duration(d)
	return
}
func (t Duration) MarshalText() (text []byte, err error) {
	return []byte(time.Duration(t).String()), nil
}

type DatasourceHTTP struct {
	URL         Template            `json:"url" yaml:"url"`
	Method      Template            `json:"method" yaml:"method"`
	QueryParams map[string]Template `json:"query_params" yaml:"query_params"`
	Headers     map[string]Template `json:"headers" yaml:"headers"`
	Body        Template            `json:"body" yaml:"body"`
	Timeout     *Duration           `json:"timeout" yaml:"timeout"`
	CacheKey    Template            `json:"cache-key" yaml:"cache-key"`
	TTL         Template            `json:"ttl" yaml:"ttl"`

	hc *http.Client

	rcache *ristretto.Cache
}

type configDatasourceHTTP struct {
	URL         string            `json:"url" yaml:"url"`
	Method      string            `json:"method" yaml:"method"`
	QueryParams map[string]string `json:"query_params" yaml:"query_params"`
	Headers     map[string]string `json:"headers" yaml:"headers"`
	Body        string            `json:"body" yaml:"body"`
	Timeout     *string           `json:"timeout" yaml:"timeout"`
	CacheKey    string            `json:"cache-key" yaml:"cache-key"`
	TTL         string            `json:"ttl" yaml:"ttl"`
}

func (dh *DatasourceHTTP) Init() (err error) {
	t := 10 * time.Second
	if dh.Timeout != nil {
		t = time.Duration(*dh.Timeout)
	}
	dh.hc = &http.Client{Timeout: t}

	dh.rcache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	return
}

func (dh *DatasourceHTTP) ToData(kong *pdk.PDK, data Data) (d map[string]interface{}, err error) {
	d, cacheKey, err := dh.getCache(data)
	if d != nil || err != nil {
		return
	}

	ur, err := dh.URL.ToString(data)
	if err != nil {
		return
	}
	u, err := url.Parse(ur)
	if err != nil {
		return
	}

	m, err := dh.Method.ToString(data)
	if err != nil {
		return
	}
	if m == "" {
		m = http.MethodGet
	}

	q := u.Query()
	for k, v := range dh.QueryParams {
		s, err := v.ToString(data)
		if err != nil {
			return nil, err
		}
		q.Add(k, s)
	}
	u.RawQuery = q.Encode()

	b, err := dh.Body.ToString(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest(m, u.String(), strings.NewReader(b))
	if err != nil {
		return
	}
	for h, t := range dh.Headers {
		v, err := t.ToString(data)
		if err != nil {
			return nil, err
		}
		req.Header.Add(h, v)
	}

	res, err := dh.hc.Do(req)
	if err != nil {
		return
	}

	h := map[string]interface{}{}
	for k, v := range res.Header {
		h[strings.ReplaceAll(strings.ToLower(k), "-", "_")] = v
	}
	rs, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var rj interface{} = map[string]interface{}{}
	if err := json.Unmarshal(rs, &rj); err != nil {
		rj = string(rs)
	}

	d = map[string]interface{}{}
	d["response"] = map[string]interface{}{
		"status_code": res.StatusCode,
		"headers":     h,
		"body":        rj,
	}
	err = dh.cache(cacheKey, d, data)

	return
}

func (dh *DatasourceHTTP) getCache(data Data) (d map[string]interface{}, key string, err error) {
	if dh.TTL.Empty() {
		return
	}

	key, err = dh.CacheKey.ToString(data)
	if err != nil {
		return
	}
	v, ok := dh.rcache.Get(key)
	if !ok {
		return
	}
	if d, ok = v.(map[string]interface{}); !ok {
		return nil, "", nil
	}
	return
}

func (dh *DatasourceHTTP) cache(key string, d map[string]interface{}, data map[string]interface{}) (err error) {
	if dh.TTL.Empty() {
		return
	}

	m := map[string]interface{}{"__self__": d}
	for k, v := range data {
		m[k] = v
	}
	s, err := dh.TTL.ToString(m)
	if err != nil {
		return
	}

	ttl, err := time.ParseDuration(s)
	if err != nil {
		d, err := cast.StringToDate(s)
		if err != nil {
			return err
		}

		now := time.Now()
		if d.Before(now) {
			ttl = 0
		} else {
			ttl = d.Sub(now)
		}
	}
	err = nil
	if ttl <= 0 {
		return
	}

	dh.rcache.SetWithTTL(key, d, 1, ttl)
	return
}
