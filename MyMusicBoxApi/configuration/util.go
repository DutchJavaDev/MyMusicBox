package configuration

import (
	"flag"
	"fmt"
	"musicboxapi/logging"
	"musicboxapi/models"
	"os"
)

var Config models.Config

func LoadConfiguration() {
	flag.StringVar(&Config.DevPort, "port", "", "-port=8081")
	flag.BoolVar(&Config.UseDevUrl, "devurl", false, "-devurl")
	flag.BoolVar(&Config.UsePlayUrl, "usePlayUrl", false, "-usePlayUrl")
	flag.BoolVar(&Config.UseImageUrl, "useImageUrl", false, "-useImageUrl")
	flag.StringVar(&Config.SourceFolder, "sourceFolder", "music", "-sourceFolder=/path to source folder/")
	flag.StringVar(&Config.OutputExtension, "outputExtension", "opus", "-outputExtension=opus,mp3,mp4 etc")
	flag.Parse()
}

func GetApiGroupUrl(version string) string {
	if Config.UseDevUrl {
		return fmt.Sprintf("/dev/api/%s", version)
	} else {
		return fmt.Sprintf("/api/%s", version)
	}
}

func DeleteFile(path string) {
	err := os.Remove(path)

	if err != nil {
		logging.ErrorStackTrace(err)
		logging.Error(fmt.Sprintf("Failed to delete file, path: %s", path))
	}
}
