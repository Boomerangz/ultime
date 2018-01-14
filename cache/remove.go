package cache

//function to remove key located in cache
func Remove(key string) error {
	return CacheInstance.Remove(key)
}
