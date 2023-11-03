package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Kong/go-pdk/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestPlugin(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"headers":         r.Header,
			"queryString":     r.URL.Query(),
			"startedDateTime": time.Now().Format(time.RFC3339Nano),
		})
	}))
	s := os.Getenv("TEST_SERVER_URL")
	os.Setenv("TEST_SERVER_URL", svr.URL)
	defer os.Setenv("TEST_SERVER_URL", s)

	defer svr.Close()
	env, err := test.New(t, test.Request{
		Method: "GET",
		Url:    "http://example.com?x_foo=bar&expire_at=" + url.QueryEscape(time.Now().Add(time.Hour).Format(time.RFC3339Nano)),
		Headers: http.Header{
			"x-foo": []string{"bar"},
		},
	})
	assert.NoError(t, err)

	f, err := os.Open("testdata/config.yaml")
	assert.NoError(t, err)

	config := &Config{}
	err = yaml.NewDecoder(f).Decode(config)
	assert.NoError(t, err)

	var mockbinDate string
	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("loop-%d", i), func(t *testing.T) {
			env.DoHttps(config)
			t.Log(env.ServiceReq.Url)
			t.Log(string(env.ClientRes.Body))

			assert.Equal(t, 200, env.ClientRes.Status)
			assert.Equal(t, 200, env.ClientRes.Status)
			assert.Equal(t, "Jon Doe", env.ServiceReq.Headers.Get("username"))
			assert.Equal(t, "world", env.ServiceReq.Headers.Get("hello"))
			assert.Equal(t, "bar", env.ServiceReq.Headers.Get("foo"))
			assert.Equal(t, "bar", env.ServiceReq.Headers.Get("foo-ex"))
			date := env.ServiceReq.Headers.Get("date")
			if mockbinDate == "" {
				mockbinDate = date
			} else {
				assert.Equal(t, mockbinDate, date)
			}
			t.Log(date)

			u, err := url.Parse(env.ServiceReq.Url)
			assert.NoError(t, err)
			assert.Equal(t, "Jon Doe", u.Query().Get("username"))
			assert.Equal(t, "world", u.Query().Get("hello"))
			assert.Equal(t, "bar", u.Query().Get("foo"))
			assert.Equal(t, "bar", u.Query().Get("foo-ex"))
			assert.Equal(t, date, u.Query().Get("date"))
		})
	}
}
