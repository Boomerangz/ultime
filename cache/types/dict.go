package types

import (
	"encoding/json"
	"sync"
)

type CacheDict struct {
	Dict  map[string]interface{}
	mutex *sync.RWMutex
}

func NewCacheDict(dict map[string]interface{}) CacheDict {
	return CacheDict{mutex: &sync.RWMutex{}, Dict: dict}
}

func (c CacheDict) GetByKey(key string) (interface{}, bool) {
	c.mutex.RLock()
	value, ok := c.Dict[key]
	c.mutex.RUnlock()
	return value, ok
}

func (c *CacheDict) SetByKey(key string, value interface{}) {
	c.mutex.Lock()
	c.Dict[key] = value
	c.mutex.Unlock()
}

func (c CacheDict) GetValue() interface{} {
	return c.Dict
}

func (c CacheDict) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.GetValue())
}
