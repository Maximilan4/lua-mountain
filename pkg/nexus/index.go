package nexus

import (
	"sync"

	"golang.org/x/exp/maps"
)

type (
	AssetIndex struct {
		mut    *sync.RWMutex
		assets map[string]Asset
	}
)

func (ai *AssetIndex) Count() int {
	ai.mut.RLock()
	defer ai.mut.RUnlock()

	return len(ai.assets)
}

// Delete - deletes an Asset from index
func (ai *AssetIndex) Delete(key string) {
	ai.mut.Lock()
	defer ai.mut.Unlock()
	delete(ai.assets, key)
}

// Store - stores an Asset in index, will rewrite old version
func (ai *AssetIndex) Store(key string, a Asset) {
	ai.mut.Lock()
	defer ai.mut.Unlock()
	ai.assets[key] = a
}

func (ai *AssetIndex) Replace(m map[string]Asset) {
	ai.mut.Lock()
	defer ai.mut.Unlock()
	ai.assets = m
}

// Keys - returns a list with all stored assets keys
func (ai *AssetIndex) Keys() []string {
	ai.mut.RLock()
	defer ai.mut.RUnlock()

	return maps.Keys(ai.assets)
}

// Get - get an asset Data from index
func (ai *AssetIndex) Get(key string) *Asset {
	ai.mut.RLock()
	defer ai.mut.RUnlock()

	if a, ok := ai.assets[key]; ok {
		return &a
	}

	return nil
}

// Has - return an existence of Asset in index
func (ai *AssetIndex) Has(key string) (ok bool) {
	ai.mut.RLock()
	defer ai.mut.RUnlock()

	_, ok = ai.assets[key]

	return
}

func NewAssetIndex() *AssetIndex {
	return &AssetIndex{
		mut:    &sync.RWMutex{},
		assets: make(map[string]Asset, 10),
	}
}
