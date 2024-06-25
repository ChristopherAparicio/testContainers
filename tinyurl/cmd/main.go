package main

import (
	"github.com/christapa/tinyurl/app"
	"github.com/christapa/tinyurl/config"
)

func main() {
	config := config.ParseConfig()
	app.Run(config)
}
