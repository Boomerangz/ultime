package cache

func Remove(key string) error {
	return CacheInstance.Remove(key)
}
