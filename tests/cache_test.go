package cache_test

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Boomerangz/ultime/cache"
	"github.com/Boomerangz/ultime/cmd/config"
)

func init() {
	cache.Init(config.Config{DataPath: "./test_files/"})
}
func TestCacheForNonExistingKey(t *testing.T) {
	actualResult, err := cache.Get("SomeNotExistingKey")
	if actualResult != nil {
		t.Fatalf("Expected %s as result but got %s", nil, actualResult)
	}
	if err == nil {
		t.Fatalf("Expected error not to be nil")
	}
}

func TestCacheForSetGet(t *testing.T) {
	key := "Key"
	settedValue := "VALUE"
	cache.Set("Key", settedValue, 0)

	actualResult, err := cache.Get(key)
	if actualResult != settedValue {
		t.Fatalf("Expected %s as result but got %s", settedValue, actualResult)
	}
	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err)
	}
}

func TestCacheForSetGetWithTTL(t *testing.T) {
	key := "Key"
	settedValue := "VALUE"
	cache.Set("Key", settedValue, 3)

	actualResult, err := cache.Get(key)
	if actualResult != settedValue {
		t.Fatalf("Expected %s as result but got %s", settedValue, actualResult)
	}
	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err)
	}
	time.Sleep(4 * time.Second)

	actualResult, err = cache.Get(key)
	if actualResult == settedValue {
		t.Fatalf("Expected result to be empty but got %v", actualResult)
	}
	if err == nil {
		t.Fatalf("Expected error to be NotFound or Expired but got %s", err)
	}

}

func TestCacheForSetGetByKey(t *testing.T) {
	key := "Key"
	internalKey := "Internal Key"
	internalValue := "Internal Value"
	newInternalValue := "New Internal Value"
	newInternalKey := "New Internal Key"
	settedValue := map[string]interface{}{
		internalKey: internalValue,
	}

	cache.Set(key, settedValue, 0)
	testForGetStructure(t, key, settedValue)
	testForGetByOldKey(t, key, internalKey, internalValue)
	testForSetByKey(t, key, internalKey, newInternalValue)
	testForSetByKey(t, key, newInternalKey, newInternalValue)
}

func TestCacheForRemove(t *testing.T) {
	key := "Key"
	value := "Value"
	cache.Set(key, value, 0)
	actualResult, err := cache.Get(key)
	if actualResult != value {
		t.Fatalf("Expected %s as result but got %s", value, actualResult)
	}
	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err)
	}
	cache.Remove(key)

	actualResult, err = cache.Get(key)
	if actualResult == value {
		t.Fatalf("Expected %v as result but got %s", nil, actualResult)
	}
	if err == nil {
		t.Fatalf("Expected error to be Not found but got %s", err)
	}
}

func TestCacheForKeys(t *testing.T) {
	key1 := "Key1"
	key2 := "Key2"
	key3 := "Completely Different String"
	value := "Value"
	cache.Set(key1, value, 0)
	cache.Set(key2, value, 0)
	cache.Set(key3, value, 0)

	actualKeys, err := cache.Keys("Key.*")
	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err)
	}
	if len(actualKeys) != 2 {
		t.Fatalf("Expected len of result to be 2 but got %s", actualKeys)
	}
}

func testForGetStructure(t *testing.T, key string, structValue map[string]interface{}) {
	actualResult, err := cache.Get(key)
	if castedResult, ok := actualResult.(map[string]interface{}); ok {
		if !reflect.DeepEqual(castedResult, structValue) {
			t.Fatalf("Expected result to be %v but got %v", structValue, castedResult)
		}
	} else {
		t.Fatalf("Expected result to be map[string]interface{} but got %s with value %s", reflect.TypeOf(castedResult), castedResult)
	}
	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err)
	}
}

func testForGetByOldKey(t *testing.T, key string, internalKey string, expectedValue string) {
	if gotValue, err := cache.GetByKey(key, internalKey); err == nil {
		if !reflect.DeepEqual(gotValue, expectedValue) {
			t.Fatalf("Expected get by key result to be %s but got %s", expectedValue, gotValue)
		}
	} else {
		t.Fatalf("Cache Get by key error %s", err)
	}
}

func testForSetByKey(t *testing.T, key string, internalKey string, internalValue string) {
	err := cache.SetByKey(key, internalKey, internalValue)
	if err != nil {
		t.Fatalf("Set by key error %s", err)
	}

	if gotValue, err := cache.GetByKey(key, internalKey); err == nil {
		if !reflect.DeepEqual(gotValue, internalValue) {
			t.Fatalf("Expected get by key result to be %s but got %s", internalValue, gotValue)
		}
	} else {
		t.Fatalf("Cache Get by key error %s", err)
	}
}

func TestPersistance(t *testing.T) {
	os.RemoveAll("./test_files/")
	os.Mkdir("./test_files/", 0777)

	key := "Key"
	settedValue := []interface{}{[]string{"1", "2"}, "3"}
	cache.Set("Key", settedValue, 0)
	log.Println("Starting save to disk")

	err := cache.CacheInstance.SaveToDisk()
	if err != nil {
		t.Fatalf("Error saving to disk %v", err)
	}

	//Remove all data from cache to test
	keys, _ := cache.CacheInstance.Keys("")
	for _, key := range keys {
		cache.CacheInstance.Remove(key)
	}

	log.Println("Starting load to disk")
	err = cache.CacheInstance.LoadFromDisk()

	if err != nil {
		t.Fatalf("Error loading from disk %v", err)
	}
	actualResult, err := cache.Get(key)
	if !reflect.DeepEqual(actualResult, settedValue) {
		t.Fatalf("Expected %s as result but got %s", settedValue, actualResult)
	}
	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err)
	}
	os.RemoveAll("./test_files")
}
