package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_createCacheFile(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	// input ""
	directory := ""
	cacheFile, err := createCacheFile(directory)
	if !assert.Equal(t, cacheFile, homeDir+"/.ArtisanCloud/cache") {
		t.Error(err)
	}

	// input "~/"
	directory = homeDir
	cacheFile, err = createCacheFile(directory)
	if !assert.Equal(t, cacheFile, homeDir+"/.ArtisanCloud/cache") {
		t.Error(err)
	}

	// input "~/test/"
	directory = homeDir + "test/"
	cacheFile, err = createCacheFile(directory)
	if !assert.Equal(t, cacheFile, homeDir+"/.ArtisanCloud/cache") {
		t.Error(err)
	}

}
