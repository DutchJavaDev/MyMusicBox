package util

import (
	"api/logging"
	"api/models"
	"flag"
)

func GetConfig() models.Config {
	var devPort string
	flag.StringVar(&devPort, "port", "", "development port")

	var devUrl bool
	flag.BoolVar(&devUrl, "devurl", false, "rewrite  modify urls for dev")

	flag.Parse()

	// only ouput my logs when in debug mode
	logging.OutputLog = devUrl

	return models.Config{
		UseDevUrl: devUrl,
		DevPort:   devPort,
	}
}

func GetApiGroupUrlV1(useDevUrl bool) string {
	if useDevUrl {
		return "/dev/api/v1"
	} else {
		return "/api/v1"
	}
}
