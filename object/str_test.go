package object

import (
	"github.com/ArtisanCloud/PowerLibs/v2/fmt"
	"testing"
)

func Test_QuickRandom(t *testing.T) {

	for i := 1; i < 5; i++ {
		response := QuickRandom(4)
		fmt.Dump(response)
	}
}
