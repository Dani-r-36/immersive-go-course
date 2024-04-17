package concurrency

import (
	"sync"
	"time"
)

//https://levelup.gitconnected.com/implementing-cache-in-golang-a8a7e631ca85

type Cache [K comparable, V any] struct {
	mu					sync.Mutex
	//need keys dict pointing to entries
	entries				map[K]*cachEntry[V]
	unsucessreads		uint64
	entryLimit			int
	removed				uint64
	removeucessread		uint64
	removeunsucessread	uint64
}
//need K as it is key pointing to value
type cachEntry[V any] struct{
	value 	V
	lastA	time.Time
	reads	uint64
}

//*Cache and not Cache as slice 
func NewCache[K comparable, V any](entryLimit int) *Cache[K, V] { 
	return &Cache[K, V]{
		entryLimit: entryLimit,
		entries: make(map[K]*cachEntry[V]),
	}
}
//remove longest staying entry and also add to the stats 
func (c *Cache[K, V]) remove() {
	if len(c.entries) == 0 {
		return
	}
	var oldentrytime time.Time
	var oldentry K
	firstEntry := true
	for k,v := range c.entries{
		//need first entry as need a reference old entry to be compared to 
		if v.lastA.Before(oldentrytime) || firstEntry{
			oldentry = k
			oldentrytime = v.lastA
		} 
	}
	removing := c.entries[oldentry]
	if removing.reads == 0{
		c.removeunsucessread++
	} else{
		c.removeucessread++
	}
	c.removed++
	delete(c.entries, oldentry)
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	//makes the location locked till exit function
	c.mu.Lock()
 	defer c.mu.Unlock()

	if len(c.entries) == c.entryLimit {
		c.remove()
	}
	//check if key already exist, with returning the entry, and by default go returns if it exists or not in found(T or F)
	entry, found := c.entries[key]
	if !found{
		//creates new entry and has to be in this exact format app
		entry = &cachEntry[V]{
			value:      value,
			lastA: time.Now(),
		}
	} else{
		//updates existing
		entry.value = value
		entry.lastA = time.Now()
	}
	//adds new or updates existing entry to entries
	c.entries[key] = entry
	return found
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	//makes the location locked till exit function
	c.mu.Lock()
 	defer c.mu.Unlock()
	//check if key already exist, with returning the entry, and by default go returns if it exists or not in found(T or F)
	entry, found := c.entries[key]
	if !found{
		c.unsucessreads++
		return nil, found
	}else{
		entry.lastA = time.Now()
		entry.reads++
		return &entry.value, found
	}
}