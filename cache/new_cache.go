package cache

import "log"

type Cache struct {
	Memory map[string][]byte
	DB     *BoltDB
}

func NewCache() *Cache {

	db := NewBoltDB()
	all := db.GetALL()

	return &Cache{
		Memory: all,
		DB:     db,
	}
}

// Get extract full URL from cache
func (c *Cache) Get(key string) []byte {
	if fullURL, ok := c.Memory[key]; !ok {
		return nil
	} else {
		return fullURL
	}
}

// Add full URL to cache using short URL as key
func (c *Cache) Add(key, fullURL string) {
	log.Printf("Added key: %s, URL: %s", key, fullURL)

	c.Memory[key] = []byte(fullURL)
	err := c.SaveToDisk(key, fullURL)
	if err != nil {
		log.Fatalf("can't save key: %s, val: %s to disk", key, fullURL)
	}
}

// SaveToDisk saves cache to disk using BoltDB
func (c *Cache) SaveToDisk(key, fullURL string) error {
	return c.DB.Put(key, fullURL)
}

// LoadFromDisk load cache from BoltDB to memory
func (c *Cache) LoadFromDisk() map[string][]byte {

	return c.DB.GetALL()
}
