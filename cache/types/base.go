package types

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrorExpiredValue error = errors.New("Value expired")

type CacheValueInterface interface {
	GetValue() interface{}
}

type CacheInt int64

func (c CacheInt) GetValue() interface{} {
	return c
}

type CacheAny struct {
	Value interface{}
}

func (c CacheAny) GetValue() interface{} {
	return c.Value
}

func (c CacheAny) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.GetValue())
}

type CacheValueWrapper struct {
	ValidUntil *time.Time
	Value      CacheValueInterface
}

func (c CacheValueWrapper) IsValid() bool {
	if c.ValidUntil != nil {
		return time.Since(*c.ValidUntil) <= 0
	} else {
		return true
	}
}

func (c CacheValueWrapper) GetValue() (CacheValueInterface, error) {
	if c.IsValid() {
		return c.Value, nil
	} else {
		value := CacheInt(0)
		return &value, ErrorExpiredValue
	}
}
