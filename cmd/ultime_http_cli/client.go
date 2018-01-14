package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Boomerangz/ultime/cmd/config"
)

var command = flag.String("command", "", "command you need to execute")
var key = flag.String("key", "", "key you need to get or set")
var internalKey = flag.String("internal_key", "", "internal key you need to get or set")
var value = flag.String("value", "", "value you need to set")
var expire = flag.Int("expire", 0, "expiration time for your value in seconds")

func main() {
	flag.Parse()
	switch *command {
	case "set":
		{
			Set(*key, *internalKey, *value, *expire)
		}
	case "get":
		{
			Get(*key, *internalKey)
		}
	case "del":
		{
			Delete(*key)
		}
	}

}

func Set(key string, internalKey string, value string, expire int) {
	if internalKey != "" {
		key = key + "/" + internalKey
	}
	var parsedValue interface{}
	json.Unmarshal([]byte(value), &parsedValue)
	code, resp, err := doPost("cache/"+key, map[string]interface{}{"value": parsedValue, "expires": expire})
	if code == 201 {
		fmt.Println("setted successfully")
	} else {
		if resp != "" {
			fmt.Println(resp)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Delete(key string) {
	code, resp, err := doDelete("cache/" + key)
	if code == 204 {
		fmt.Println("Removed succefully")
	} else {
		if resp != "" {
			fmt.Println(resp)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Get(key string, internalKey string) {
	if internalKey != "" {
		key = key + "/" + internalKey
	}
	code, resp, err := doGet("cache/" + key)
	if err != nil {
		fmt.Println(err)
		if resp != "" {
			fmt.Println(resp)
		}
	} else if code == 200 {
		var parsedMap map[string]interface{}
		err = json.Unmarshal([]byte(resp), &parsedMap)
		if err != nil {
			fmt.Println(err)
		} else {
			marshalled, _ := json.Marshal(parsedMap["value"])
			fmt.Println(string(marshalled))
		}
	} else {
		if resp != "" {
			fmt.Println(resp)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func doDelete(path string) (int, string, error) {
	url := fmt.Sprintf("http://%s/%s", config.GetConfig().ServerUrl, path)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), err
}

func doGet(path string) (int, string, error) {
	url := fmt.Sprintf("http://%s/%s", config.GetConfig().ServerUrl, path)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), err
}
func doPost(path string, message map[string]interface{}) (int, string, error) {
	url := fmt.Sprintf("http://%s/%s", config.GetConfig().ServerUrl, path)
	jsonStr, err := json.Marshal(message)
	if err != nil {
		return 0, "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), err
}
