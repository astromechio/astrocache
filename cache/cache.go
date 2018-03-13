package cache

// Cache defines the cache implementaion for astrocache
type Cache struct {
	Values map[string]string
}

// EmptyCache returns an empty cache
func EmptyCache() *Cache {
	return &Cache{
		Values: make(map[string]string),
	}
}

// SetValueForKey sets a value for a key
func (c *Cache) SetValueForKey(val, key string) {
	c.Values[key] = val
}

// ValueForKey retreives a value for a key
func (c *Cache) ValueForKey(key string) string {
	if val, ok := c.Values[key]; ok {
		return val
	}

	return ""
}
