package util

import (
	"api/models"
	"flag"
)

var Config models.Config

func LoadConfig() {
	var devPort string
	flag.StringVar(&devPort, "port", "", "development port")

	var devUrl bool
	flag.BoolVar(&devUrl, "devurl", false, "rewrite  modify urls for dev")

	flag.Parse()

	Config.DevPort = devPort
	Config.UseDevUrl = devUrl
}

func GetApiGroupUrlV1(useDevUrl bool) string {
	if useDevUrl {
		return "/dev/api/v1"
	} else {
		return "/api/v1"
	}
}
