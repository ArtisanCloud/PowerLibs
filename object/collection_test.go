package object

import (
	"github.com/ArtisanCloud/PowerLibs/fmt"
	"testing"
)

func Test_Collection_Set_AND_Get(t *testing.T) {

	collectionTest := NewCollection(&HashMap{
		"gun": "model",
	})

	collectionTest.Set("weapon.bullet", 100)
	collectionTest.Set("weapon.shield.strength", "strong")

	bulletCount := collectionTest.Get("weapon.bullet", 0)
	if bulletCount != 100 {
		t.Error("get bullet error")
		fmt.Dump(bulletCount)
	}

	shieldStrength := collectionTest.Get("weapon.shield.strength", "")
	if shieldStrength != "strong" {
		t.Error("get shield error")
		fmt.Dump(shieldStrength)
	}

}
