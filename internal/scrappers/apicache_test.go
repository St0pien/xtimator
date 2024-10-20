package scrappers

import (
	"os"
	"reflect"
	"testing"
)

func TestApiCacheBasicOperations(t *testing.T) {
	cache, err := NewApiCache(".cache", "test")
	if err != nil {
		t.Fatalf("Couldn't create cache: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(".cache")
	})

	var testBuf [1024]byte
	for i := 0; i < 1024; i++ {
		testBuf[i] = byte(i)
	}
	err = cache.Save("cached", testBuf[:])
	if err != nil {
		t.Fatalf("Failed to save %v", err)
	}

	result, err := cache.Get("cached")
	if err != nil {
		t.Fatalf("Failed to get %v", err)
	}

	if !reflect.DeepEqual(result, testBuf[:]) {
		t.Fatal("Bytes differ", result, testBuf)
	}

}
