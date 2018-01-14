package cache

func Keys(keyPattern string) ([]string, error) {
	return CacheInstance.Keys(keyPattern)
}
