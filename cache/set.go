package cache

import (
	"strconv"

	"github.com/boomerangz/ultime/cache/types"
)

//function to set value to cache, requires key, value to set.
//also function receives Expires parameter, which hadles Time to store value in seconds
//if expires == 0 then value is stored infinetely
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

//function to set value into stored structure (array or dict)
//receives cache key of structure, internal key that references to inside structure place and value to set
//it is impossible to set Expiring value inside structure
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
