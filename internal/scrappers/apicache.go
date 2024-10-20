package scrappers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ApiCache struct {
	BasePath string
	lookup   map[string]int
}

func NewApiCache(path ...string) (ApiCache, error) {
	cachePath := filepath.Join(path...)
	cache := ApiCache{BasePath: cachePath}

	lookupFilePath := filepath.Join(cachePath, "lookup.json")
	lookupFile, err := os.ReadFile(lookupFilePath)

	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(cachePath, os.ModePerm); err != nil {
			return ApiCache{}, fmt.Errorf("couldn't create cache dir: %w", err)
		}
		cache.lookup = make(map[string]int)
	} else if err != nil {
		return ApiCache{}, fmt.Errorf("failed to read lookup file: %w", err)
	} else if err := json.Unmarshal(lookupFile, &cache.lookup); err != nil {
		return ApiCache{}, fmt.Errorf("failed to parse lookup file: %w", err)
	}

	return cache, nil
}

func (cache *ApiCache) Get(key string) ([]byte, error) {
	fileId, ok := cache.lookup[key]

	if !ok {
		return nil, nil
	}

	content, err := os.ReadFile(filepath.Join(cache.BasePath, fmt.Sprintf("%v.json", fileId)))
	if err != nil {
		return nil, fmt.Errorf("Failed reading from cache [%v]: %w", fileId, err)
	}

	return content, nil
}

func (cache *ApiCache) Save(key string, content []byte) error {
	nextFileId := len(cache.lookup)

	err := os.WriteFile(filepath.Join(cache.BasePath, fmt.Sprintf("%v.json", nextFileId)), content, 0644)
	if err != nil {
		return fmt.Errorf("Failed writing to cache [%v]: %w", nextFileId, err)
	}

	cache.lookup[key] = nextFileId

	encodedLookup, err := json.Marshal(cache.lookup)
	if err != nil {
		delete(cache.lookup, key)
		return fmt.Errorf("Failed serializing lookup: %w", err)
	}

	err = os.WriteFile(filepath.Join(cache.BasePath, "lookup.json"), encodedLookup, 0644)
	if err != nil {
		delete(cache.lookup, key)
		return fmt.Errorf("Failed saving lookup: %w", err)
	}

	return nil
}
