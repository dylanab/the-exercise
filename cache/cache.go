package cache;

import ( 
	"time"
	"fmt"
	"log"
);

type CacheItem struct{
	data []byte;
	timeCached time.Time;
	expirationTime time.Time;
}

type Cache struct{
	storage map[string]*CacheItem;
	itemDuration time.Duration;
	canServeStale bool;
}

func Create(itemDuration string, canServeStale bool) (*Cache, error) {
	cache := Cache{};
	cache.storage = make(map[string]*CacheItem);
	d, err := time.ParseDuration(itemDuration);
	if(err != nil) {
		return nil, fmt.Errorf("failed to create cache: %v", err);
	}
	cache.itemDuration = d;
	cache.canServeStale = canServeStale;
	return &cache, nil;
}

func (c *Cache) AddItem(key string, data []byte) {	
	_, ok := c.storage[key];
	if(ok) {
		log.Printf("overwriting item in cache");
	}
	item  := CacheItem{};
	item.data = data;
	item.timeCached = time.Now();
	item.expirationTime = item.timeCached.Add(c.itemDuration);
	c.storage[key] = &item;
}

func (c *Cache) GetItem(key string) ([]byte, bool) {
	item, ok := c.storage[key];
	if(!ok) {
		log.Printf("item %s not found in cache", key);
		return []byte(""), true;
	}
	
	log.Printf("item %s was found in cache", key);

	data := item.data;
	now := time.Now();
	itemIsExpired := false;
	if(item.expirationTime.After(now)) {
		log.Printf("item %s is expired", key);
		itemIsExpired = true;
		if(!c.canServeStale) {
			log.Printf("item %s cannot be served stale, deleting it.", key);
			data = []byte("");
			delete(c.storage, key);
		}
	}

	return data, itemIsExpired;
}	