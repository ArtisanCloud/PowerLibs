package object

import (
	"github.com/ArtisanCloud/PowerLibs/fmt"
	"testing"
)

func Test_QuickRandom(t *testing.T) {

	response := QuickRandom(32)
	fmt.Dump(response)

}
