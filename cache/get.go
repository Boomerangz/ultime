package cache

import (
	"strconv"

	"github.com/boomerangz/ultime/cache/types"
)

//Function to get value out of cache, requires a key
func Get(key string) (interface{}, error) {
	valueWrapper, err := CacheInstance.Get(key)
	if err != nil {
		return nil, err
	}
	value, err := valueWrapper.GetValue()
	return value.GetValue(), err
}

//Function to get value out structure (array or dict) stored in cache, requires a key and an internal key
func GetByKey(key string, internalKey string) (interface{}, error) {
	valueWrapper, err := CacheInstance.Get(key)
	if err != nil {
		return nil, err
	}
	value, err := valueWrapper.GetValue()
	if err != nil {
		return nil, err
	}
	if castedMap, ok := value.(types.CacheDict); ok {
		if mapValue, ok := castedMap.GetByKey(internalKey); ok {
			if cvi, ok := mapValue.(types.CacheValueInterface); ok {
				return cvi.GetValue(), nil
			}
			return mapValue, nil
		} else {
			return nil, NotFoundError
		}
	} else if castedList, ok := value.(types.CacheList); ok {
		if intKey, err := strconv.Atoi(internalKey); err == nil {
			listValue, ok := castedList.GetByIndex(intKey)
			if !ok {
				return nil, NotFoundError
			}
			if cvi, ok := listValue.(types.CacheValueInterface); ok {
				return cvi.GetValue(), nil
			}
			return listValue, nil
		} else {
			return nil, NotFoundError
		}
	} else {
		return nil, NotFoundError
	}
}
