package configuration

import (
	"flag"
	"fmt"
	"musicboxapi/models"
)

var Config models.Config

func LoadConfig() {
	flag.StringVar(&Config.DevPort, "port", "", "Development port else use default port")
	flag.BoolVar(&Config.UseDevUrl, "devurl", false, "Have a dev prefix in the url")
	flag.StringVar(&Config.SourceFolder, "sourceFolder", "music", "Output folder for data")
	flag.StringVar(&Config.OutputExtension, "outputExtension", "opus", "Extension for ouput file")
	flag.Parse()
}

func GetApiGroupUrlV1() string {
	if Config.UseDevUrl {
		return "/dev/api/v1"
	} else {
		return "/api/v1"
	}
}

func GetApiGroupUrl(version string) string {
	if Config.UseDevUrl {
		return fmt.Sprintf("/dev/api/%s", version)
	} else {
		return fmt.Sprintf("/api/%s", version)
	}

}
