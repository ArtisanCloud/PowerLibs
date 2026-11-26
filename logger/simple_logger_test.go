package logger

import (
	"bytes"
	"log"
	"testing"

	"github.com/ArtisanCloud/PowerLibs/v3/logger/drivers"
)

func TestSimpleLoggerInfo(t *testing.T) {
	var output bytes.Buffer
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	t.Cleanup(func() {
		log.SetOutput(writer)
		log.SetFlags(flags)
		log.SetPrefix(prefix)
	})
	log.SetOutput(&output)
	log.SetFlags(0)
	log.SetPrefix("")

	logger, err := NewLogger(&drivers.SimpleLogger{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		log  func()
		want string
	}{
		{"callback", func() { logger.Info("Server response created:", "content", "12345") }, "Server response created: content=12345"},
		{"request fields", func() { logger.Info("Request received:", "method", "GET", "uri", "/a%20b") }, "Request received: method=GET uri=/a%20b"},
		{"literal message", func() { logger.Info("progress 100%") }, "progress 100%"},
		{"legacy printf", func() { logger.Info("%s/%s", "a", "b") }, "a/b"},
		{"formatted arguments", func() { logger.InfoF("%s size=%d", "/a%20b", 42) }, "/a%20b size=42"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output.Reset()
			test.log()
			if got, want := output.String(), "[INFO] "+test.want+"\n"; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
