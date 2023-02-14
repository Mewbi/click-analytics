package main

import (
	"click-analytics/config"
	"click-analytics/server"
)

func main() {
    config.Load("config/config.toml")
    server.Start()
}
