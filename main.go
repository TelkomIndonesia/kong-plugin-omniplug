package main

import (
	"os"
	"strconv"

	"github.com/Kong/go-pdk/server"
)

var version = "0.1.0"
var priority = func() int {
	s := os.Getenv("KONG_PLUGIN_OMNIPLUG_PRIORITY")
	i, err := strconv.Atoi(s)
	if err != nil {
		return 810
	}
	return i
}()

func main() {
	server.StartServer(NewConfig, version, priority)
}
