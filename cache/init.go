package cache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/Boomerangz/ultime/cmd/config"
	"github.com/boomerangz/ultime/cache/types"
)

type Cache struct {
	cacheInstance map[string]types.CacheValueWrapper
	cacheMutex    sync.RWMutex

	dataPath       string
	savingInterval uint32
}

var CacheInstance *Cache
var NotFoundError error = errors.New("No such value")

func Init(config config.Config) {
	CacheInstance = &Cache{
		cacheInstance:  make(map[string]types.CacheValueWrapper),
		dataPath:       config.DataPath,
		savingInterval: config.SavingInterval,
	}
	go CacheInstance.deleteExpiredData()
	registerTypesInGob()
	if CacheInstance.dataPath != "" {
		err := CacheInstance.LoadFromDisk()
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		} else if os.IsNotExist(err) {
			log.Printf("Data file %s does not exit. Starting with empty cache\n", CacheInstance.getDataFilePath())
		}
		go CacheInstance.backgroundSaving()
	}
}

func (c *Cache) Get(key string) (types.CacheValueWrapper, error) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()
	if value, ok := c.cacheInstance[key]; ok {
		return value, nil
	} else {
		return types.CacheValueWrapper{}, NotFoundError
	}
}

func (c *Cache) Set(key string, value types.CacheValueInterface, expiresSeconds int) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	wrapper := types.CacheValueWrapper{}
	wrapper.Value = value
	if expiresSeconds > 0 {
		validUntil := time.Now().Local().Add(time.Second * time.Duration(expiresSeconds))
		wrapper.ValidUntil = &validUntil
	}
	c.cacheInstance[key] = wrapper
}

func (c *Cache) Remove(key string) error {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	if _, ok := c.cacheInstance[key]; ok {
		delete(c.cacheInstance, key)
		return nil
	} else {
		return NotFoundError
	}
}

func (c *Cache) Keys(keyPattern string) ([]string, error) {
	var re *regexp.Regexp
	if keyPattern != "" {
		localRe, err := regexp.Compile(keyPattern)
		if err != nil {
			return nil, err
		}
		re = localRe
	}
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()
	result := []string{}
	for key, value := range c.cacheInstance {
		if (re == nil || re.MatchString(key)) && value.IsValid() {
			result = append(result, key)
		}
	}
	return result, nil
}

func (c *Cache) deleteExpiredData() {
	for {
		time.Sleep(30 * time.Second)
		for key, value := range c.cacheInstance {
			if !value.IsValid() {
				c.cacheMutex.Lock()
				if !c.cacheInstance[key].IsValid() {
					delete(c.cacheInstance, key)
					log.Printf("Cleaning %s from cache\n", key)
				}
				c.cacheMutex.Unlock()
			}
		}
	}
}

func (c *Cache) backgroundSaving() {
	if c.savingInterval <= 0 {
		return
	}
	for {
		time.Sleep(time.Duration(c.savingInterval) * time.Second)
		err := c.SaveToDisk()
		if err != nil {
			log.Printf("Saving to disk problem: %v", err)
		} else {
			log.Println("Background saving finished")
		}
	}
}

func (c *Cache) SaveToDisk() error {
	copyCache := c.cacheInstance
	err := os.Mkdir(c.dataPath, 0777)
	fmt.Println("create dir", c.dataPath)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create(c.getDataFilePath())
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(copyCache)
	if err != nil {
		return err
	}
	err = file.Close()
	return err
}

func (c *Cache) LoadFromDisk() error {
	log.Printf("Load data from %s started\n", c.getDataFilePath())
	file, err := os.Open(c.getDataFilePath())
	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&c.cacheInstance)
	} else {
		log.Printf("Open file %s finished with error: %v\n", c.getDataFilePath(), err)
		return err
	}
	if err == nil {
		log.Printf("Decoding data from %s finished successfully\n", c.getDataFilePath())
	} else {
		log.Printf("Decoding data from %s finished with error: %v\n", c.getDataFilePath(), err)
	}

	return err
}

func (c *Cache) getDataFilePath() string {
	return path.Join(c.dataPath, "data.gob")
}

func registerTypesInGob() {
	cacheInt := types.CacheInt(0)
	gob.Register(cacheInt)
	gob.Register(types.CacheDict{})
	gob.Register(types.CacheAny{})
	gob.Register(types.CacheList{})
	gob.Register(types.CacheValueWrapper{})
}
