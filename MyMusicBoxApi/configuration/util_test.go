package configuration

import (
	"flag"
	"musicboxapi/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	// Arrange
	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"cmd", "-port=8081", "-devurl", "-sourceFolder=dev_music", "-outputExtension=opus"}

	// Act
	LoadConfiguration()

	// Assert
	var config models.Config = Config

	assert.Equal(t, "8081", config.DevPort)
	assert.Equal(t, "opus", config.OutputExtension)
	assert.Equal(t, "dev_music", config.SourceFolder)
	assert.Equal(t, true, config.UseDevUrl)
}
