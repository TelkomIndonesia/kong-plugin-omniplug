package main

import (
	"net/http"
	"testing"

	"github.com/Kong/go-pdk/test"
	"github.com/alecthomas/assert/v2"
)

func TestPlugin(t *testing.T) {
	env, err := test.New(t, test.Request{
		Method: "GET",
		Url:    "http://example.com?q=search&x=9",
		Headers: http.Header{
			"X-Test": []string{"test"},
		},
	})
	assert.NoError(t, err)

	tmpl1, err := NewTemplate("{{ (index (index .data 0) (index (index .data 1).request_headers.x_test 0)).message }}")
	assert.NoError(t, err)

	re := &ReqEnrich{
		Datasource: []Datasource{
			{
				Type: "static",
				Static: map[string]interface{}{
					"test": map[string]interface{}{
						"message": "world",
					},
				},
			},
			{
				Type: "kong",
				Kong: DatasourceKong{
					RequestHeaders: true,
				},
			},
		},
		Injector: Injector{
			Headers: map[string]Template{
				"hello": tmpl1,
			},
		},
	}
	env.DoHttps(re)
	t.Log(string(env.ClientRes.Body))
	assert.Equal(t, 200, env.ClientRes.Status)
	assert.Equal(t, "world", env.ServiceReq.Headers.Get("hello"))
}
