package util

import (
	"api/logging"
	"flag"
)

func GetConfig() Config {
	var devPort string
	flag.StringVar(&devPort, "port", "", "development port")

	var devUrl bool
	flag.BoolVar(&devUrl, "devurl", false, "rewrite  modify urls for dev")

	flag.Parse()

	// only ouput my logs when in debug mode
	logging.OutputLog = devUrl

	return Config{
		UseDevUrl: devUrl,
		DevPort:   devPort,
	}
}

func GetApiGroupUrl(useDevUrl bool) string {
	if useDevUrl {
		return "/dev/api"
	} else {
		return "/api"
	}
}
