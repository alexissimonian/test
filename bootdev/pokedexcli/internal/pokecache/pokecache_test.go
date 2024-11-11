package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(time.Second)
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
		{
			key: "key2",
			val: []byte("val2"),
		},
		{
			key: "",
			val: []byte("val3"),
		},
	}

	cache := NewCache(time.Second)

	for _, cas := range cases {
		cache.Add(cas.key, cas.val)
		actual, ok := cache.Get(cas.key)
		if !ok {
			t.Errorf("Added case %s not found", cas.key)
			continue
		}

		if string(actual) != string(cas.val) {
			t.Errorf("Error getting right cache. Expected %v; Got %v",
				string(cas.val),
				string(actual),
			)
			continue
		}
	}
}

func TestReapCache(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)
	keyOne := "key1"
	cache.Add(keyOne, []byte("val1"))
	time.Sleep(interval + time.Millisecond)

	_, ok := cache.Get(keyOne)
	if ok {
		t.Errorf("%s was not removed", keyOne)
	}
}
