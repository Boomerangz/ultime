package types

import (
	"encoding/json"
	"errors"
	"sync"
)

type CacheList struct {
	List  []interface{}
	mutex sync.RWMutex
}

func NewCacheList(list []interface{}) CacheList {
	return CacheList{List: list}
}

func (c CacheList) GetLength() int {
	return len(c.List)
}

func (c *CacheList) GetByIndex(idx int) (interface{}, bool) {
	if idx >= len(c.List) {
		return nil, false
	}
	c.mutex.RLock()
	value := c.List[idx]
	c.mutex.RUnlock()
	return value, true
}

func (c *CacheList) SetByIndex(idx int, value interface{}) error {
	if idx < 0 || idx >= len(c.List) {
		return errors.New("No such item")
	}
	c.mutex.Lock()
	c.List[idx] = value
	c.mutex.Unlock()
	return nil
}

func (c CacheList) Append(idx int, value interface{}) {
	c.mutex.Lock()
	c.List = append(c.List, value)
	c.mutex.Unlock()
}

func (c CacheList) GetValue() interface{} {
	return c.List
}

func (c CacheList) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.GetValue())
}
