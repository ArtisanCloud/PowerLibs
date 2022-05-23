package gout

import (
	"github.com/ArtisanCloud/PowerLibs/object"
	"github.com/guonaihong/gout"
)

func ConvertHashMapToGoutMap(hashMap *object.HashMap) gout.H{
	gMap:=gout.H{}
	for k,v :=range *hashMap{
		gMap[k] = v
	}
	return gMap
}