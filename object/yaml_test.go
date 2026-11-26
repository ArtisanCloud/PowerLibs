package object_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"go.yaml.in/yaml/v3"
)

type uppercaseLabel string

func (label *uppercaseLabel) UnmarshalYAML(node *yaml.Node) error {
	var value string
	if err := node.Decode(&value); err != nil {
		return err
	}
	*label = uppercaseLabel(strings.ToUpper(value))
	return nil
}

func TestYAMLFileRoundTrip(t *testing.T) {
	type settings struct {
		Name  string         `yaml:"name"`
		Label uppercaseLabel `yaml:"label"`
	}
	path := filepath.Join(t.TempDir(), "settings.yaml")
	if err := object.SaveYMLFile(settings{Name: "worker", Label: "ready"}, path, 0600); err != nil {
		t.Fatal(err)
	}
	var got settings
	if err := object.OpenYMLFile(path, &got); err != nil {
		t.Fatal(err)
	}
	want := settings{Name: "worker", Label: "READY"}
	if got != want {
		t.Errorf("decoded settings = %+v, want %+v", got, want)
	}
}
