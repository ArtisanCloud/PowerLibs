package zap

import (
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"net/http"
	"testing"
)

var strArtisanCloudPath = "/var/log/ArtisanCloud/PowerLibs"
var strOutputPath = strArtisanCloudPath + "/output.log"
var strErrorPath = strArtisanCloudPath + "/errors.log"

func Test_Log_Info(t *testing.T) {
	logger, err := NewLogger(&object.HashMap{
		"env":        "test",
		"outputPath": strOutputPath,
		"errorPath":  strErrorPath,
	})
	if err != nil {
		t.Error(err)
	}

	logger.Info("test info", "response", &http.Response{})

}
