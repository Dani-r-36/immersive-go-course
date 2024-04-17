package concurrency

import (
	"reflect"
	"testing"
	"time"
)


func TestPutGet(t *testing.T){
	t.Run("empty get", func(t *testing.T) {
		cache := NewCache[string, string](1)
		_, exists := cache.Get("random")
		reflect.DeepEqual(exists, false)
	})
	t.Run("put success and get value", func(t *testing.T) {
		cache := NewCache[string, string](1)
		exists := cache.Put("key", "value")
		reflect.DeepEqual(exists, false)
		_, exists = cache.Get("key")
		reflect.DeepEqual(exists, true)
	})
	t.Run("put twice and check value and time", func(t *testing.T) {
		cache := NewCache[string, string](2)
		exists := cache.Put("key", "value")
		reflect.DeepEqual(exists, false)
		exists = cache.Put("key", "new value")
		entry, _ := cache.Get("key")
		reflect.DeepEqual("new value", *entry)
		roundedParsedTime := cache.entries["key"].lastA.Truncate(time.Second)
		currentTime := time.Now()
		roundedCurrentTime := currentTime.Truncate(time.Second)
		reflect.DeepEqual(roundedParsedTime, roundedCurrentTime)
	})


	}

func TestReads(t *testing.T){
	t.Run("check unsucessreads and reads", func(t *testing.T) {
		cache := NewCache[string, string](1)
		_ = cache.Put("key", "value")
		_, _ = cache.Get("random")
		reflect.DeepEqual(cache.entries["key"].reads, 1)
		newCacheCheck := NewCache[string, string](1)
		_, _ = newCacheCheck.Get("random")
		reflect.DeepEqual(newCacheCheck.unsucessreads, 1)
	})
	}

func TestRemove(t *testing.T){
	t.Run("check unsucessreads and reads", func(t *testing.T) {
		cache := NewCache[string, string](2)
		_ = cache.Put("key", "value")
		_, _ = cache.Get("random")
		_ = cache.Put("newkey", "new")
		_ = cache.Put("last", "new")
		reflect.DeepEqual(cache.removeucessread, 1)
		reflect.DeepEqual(cache.removed, 1)
		_, exists := cache.Get("key")
		reflect.DeepEqual(exists, false)
	})
	t.Run("check unsucessreads and reads", func(t *testing.T) {
		cache := NewCache[string, string](2)
		_ = cache.Put("key", "value")
		_, _ = cache.Get("random")
		_ = cache.Put("newkey", "new")
		_ = cache.Put("last", "new")
		reflect.DeepEqual(cache.removeucessread, 1)
		reflect.DeepEqual(cache.removed, 1)
		_, exists := cache.Get("key")
		reflect.DeepEqual(exists, false)
		_ = cache.Put("unsucess", "new")
		reflect.DeepEqual(cache.removeucessread, 0)
		reflect.DeepEqual(cache.removeunsucessread, 1)
		reflect.DeepEqual(cache.removed, 2)
	})
	}