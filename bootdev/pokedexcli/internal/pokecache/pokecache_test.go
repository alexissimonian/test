package pokecache

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()
	if cache.cache == nil {
		t.Error("Error creating new cache")
	}
}

func TestAddGetCache(t *testing.T) {
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "key1",
			val: []byte("val1"),
		},
	}

	cache := NewCache()

	for _, cas := range cases {
		cache.Add(cas.key, cas.val)
		actual, ok := cache.Get(cas.key)
		if !ok {
			t.Error("Added case not found")
		}

		if string(actual) != string(cas.val) {
			t.Errorf("Error getting right cache. Expected %v; Got %v", string(cas.val), string(actual))
		}
	}
}
