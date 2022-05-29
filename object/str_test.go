package object

import (
	"github.com/ArtisanCloud/PowerLibs/v2/fmt"
	"testing"
)

func Test_QuickRandom(t *testing.T) {

	response := QuickRandom(32)
	fmt.Dump(response)

}
