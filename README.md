# ULTIME

## Your 2K18 Ultimate Golang-powered alternative to REDIS

#### *(Not really)*

### Install:
`go get github.com/boomerangz/ultime/cmd/ultime_http`
### Run:
`ultime_http`

### Server environment variables

If you want to tune server's behavior, you can set environment variables
`ULTIME_DATAPATH` that contains a folder for server to store it's data  
Default value is `/tmp/ultime/`  
`ULTIME_SAVINGINTERVAL` that contains a interval in seconds for server to save it's data to disk. If it is set to 0, server will save it's data only at shutdown.  
Default value is `300`.  
`ULTIME_PORT` that contains a port to listen.  
Default value is `8080`


HTTP servers listens to your requests in simple format:
`http://<hostname>/cache/<key>/`

### Get a value
for example, to get value `A`:  
`curl http://127.0.0.1:8080/cache/A/`


### Set a value
for setting value `A` to `Some String Value` you just need to send POST request to the same path with JSON-encoded body.
JSON encoded body must contain field `value` with value you want to set,
and *CAN* contain not necessary field `expires` that contains field with Seconds.
If there is no value in `expires` field or it equals to 0, so there will be no expiration time for this value.
 
```curl -X POST
  http://127.0.0.1:8080/cache/a 
  -H 'Content-Type: application/json' 
  -d '{  
	"value":"asdasdadasdasd",
	"expires": 30
}'
```

### Set a structure
Also you can insert a structure (array or dictionary) into cache. You can do it simple sending it in JSON encoded body. For example:

```curl -X POST  
  http://127.0.0.1:8080/cache/a 
  -H 'Content-Type: application/json' 
  -d '{  
	"value":["asdasdadasdasd", 123123]
}'
```

### Get internal value of a structure by key 
And you also can access to one single value of your structure with putting it's internal key into path, like:
```
curl http://127.0.0.1:8080/cache/a/0
```
or 
```
curl http://127.0.0.1:8080/cache/a/1
```

### Set internal value of a structure by key 
Setting operation is also permitted to that internal value:
```curl -X POST  
  http://127.0.0.1:8080/cache/a/0 
  -H 'Content-Type: application/json' 
  -d '{  
	"value":0
}' 
```
### Remove value from cache
```
curl -X DELETE http://127.0.0.1:8080/cache/a
```

### Get cache keys list
You can receive a list of keys presented in cache.
Also you can filter it by some regular expression

```
curl http://127.0.0.1:8080/cache-keys[?pattern=<regexp>]
```






