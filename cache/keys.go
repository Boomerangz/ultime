package cache

//Function to get all keys located in cache
func Keys(keyPattern string) ([]string, error) {
	return CacheInstance.Keys(keyPattern)
}
