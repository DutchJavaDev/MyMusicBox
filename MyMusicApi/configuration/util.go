package configuration

import (
	"flag"
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

func GetApiGroupUrlV1(useDevUrl bool) string {
	if useDevUrl {
		return "/dev/api/v1"
	} else {
		return "/api/v1"
	}
}
