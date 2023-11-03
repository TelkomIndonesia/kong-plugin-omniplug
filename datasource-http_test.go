package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustTemplate(t *testing.T, s string) Template {
	tmpl, err := NewTemplate(s)
	require.NoErrorf(t, err, "invalid template %s", s)
	return tmpl
}

func TestHTTP(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"expire_at": time.Now().Add(3 * time.Second).Format(time.RFC3339Nano),
		})
	}))

	ds := &DatasourceHTTP{
		URL: mustTemplate(t, svr.URL),
		TTL: mustTemplate(t, "{{.__self__.response.body.expire_at}}"),
	}
	err := ds.Init()
	require.NoErrorf(t, err, "error initiating datasource: %s", err)

	_, err = ds.ToData(nil, map[string]interface{}{})
	require.NoErrorf(t, err, "error generating data: %s", err)

	<-time.After(time.Second)
	_, ok := ds.rcache.Get("")
	assert.True(t, ok, "response should be cached")

	<-time.After(3 * time.Second)
	_, ok = ds.rcache.Get("")
	assert.False(t, ok, "cache should no longer exist")
}
