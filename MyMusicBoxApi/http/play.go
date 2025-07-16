package http

import (
	"fmt"
	"musicboxapi/configuration"
	"musicboxapi/logging"
	"musicboxapi/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Play(ctx *gin.Context) {
	sourceId := ctx.Param("sourceId")

	if sourceId == "" {
		ctx.Status(http.StatusInternalServerError)
		logging.Warning("No sourceId in request")
		return
	}

	path := fmt.Sprintf("%s/%s.%s", configuration.Config.SourceFolder, sourceId, configuration.Config.OutputExtension)

	file, err := os.Open(path)

	if err != nil {
		logging.Error(fmt.Sprintf("Failed to read file, %s", err.Error()))
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logging.Error(fmt.Sprintf("Failed could not get fileinfo: %s", err.Error()))
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	// Set proper headers for OPUS audio
	ctx.Header("Content-Type", "audio/opus")
	ctx.Header("Accept-Ranges", "bytes")

	// Use http.ServeContent which handles range requests automatically
	http.ServeContent(
		ctx.Writer,
		ctx.Request,
		fileInfo.Name(),
		fileInfo.ModTime(),
		file,
	)
}
