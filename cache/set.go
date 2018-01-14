package cache

import (
	"strconv"

	"github.com/boomerangz/ultime/cache/types"
)

func Set(key string, value interface{}, expires int) error {
	CacheInstance.Set(key, wrapValue(value), expires)
	return nil
}

func wrapValue(value interface{}) types.CacheValueInterface {
	if casted, ok := value.(int); ok {
		result := types.CacheInt(casted)
		return result
	} else if casted, ok := value.(map[string]interface{}); ok {
		result := types.NewCacheDict(casted)
		return result
	} else if casted, ok := value.([]interface{}); ok {
		result := types.NewCacheList(casted)
		return result
	} else {
		result := types.CacheAny{Value: value}
		return result
	}
}

func SetByKey(key string, internalKey string, setValue interface{}) error {
	valueWrapper, err := CacheInstance.Get(key)
	if err != nil {
		return err
	}
	value, err := valueWrapper.GetValue()
	if err != nil {
		return err
	}
	if castedMap, ok := value.(types.CacheDict); ok {
		castedMap.SetByKey(internalKey, setValue)
		return nil
	} else if castedList, ok := value.(types.CacheList); ok {
		if intKey, err := strconv.Atoi(internalKey); err == nil {
			err := castedList.SetByIndex(intKey, setValue)
			return err
		} else {
			return NotFoundError
		}
	} else {
		return NotFoundError
	}
}
