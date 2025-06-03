package util

import (
	"api/models"
	"bufio"
	"flag"
	"os"
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

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
