package main

import (
	"MailingBot/internal/app"
	"flag"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config/config.yml", "Path to configuration file")
}

func main() {
	app.Run(configFile)
}
