package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ArtisanCloud/PowerLibs/v3/logger/config"
)

func TestLoadConfigFromFile(t *testing.T) {
	want := &config.LogConfig{
		Level:   "debug",
		Console: true,
		File:    &config.FileConfig{Enable: true, FilePath: "app.log"},
	}
	tests := []struct {
		name string
		data string
	}{
		{"yaml", "level: debug\nconsole: yes\nfile:\n  enable: true\n  filePath: app.log\n"},
		{"json", `{"level":"debug","console":true,"file":{"enable":true,"filePath":"app.log"}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "logging."+test.name)
			if err := os.WriteFile(path, []byte(test.data), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := config.LoadConfigFromFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("decoded config = %+v, want %+v", got, want)
			}
		})
	}
}
