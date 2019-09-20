package cache

import "log"

type Cache struct {
	Memory map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		Memory: map[string][]byte{},
	}
}

func (c *Cache) Get(key string) []byte {
	if fullURL, ok := c.Memory[key]; !ok {
		return nil
	} else {
		return fullURL
	}
}

func (c *Cache) Add(key, fullURL string) {
	log.Printf("Added key: %s, URL: %s", key, fullURL)

	c.Memory[key] = []byte(fullURL)
}
