# ULTIME CLIENT

## Your 2K18 Ultimate Golang-powered alternative to REDIS client



### Install:
`go get github.com/boomerangz/ultime/cmd/ultime_http_cli`
### Run:
`ultime_http_cli <params>`

### Client environment variables

If you want to tune client's behavior, you can set environment variables  
`ULTIME_SERVERURL` that contains a server hostname  
Default value is `127.0.0.1:8080`.


### Get a value
`ultime_http_cli -command get -key A`


### Set a value
`ultime_http_cli -command set -key A -value B`

### Set a structure
`ultime_http_cli -command set -key A -value '{"B":"C"}'`

### Get internal value of a structure by key 
`ultime_http_cli -command get -key A -internal_key B`

### Set internal value of a structure by key 
`ultime_http_cli -command set -key A -internal_key B -value C`


### Remove value from cache
`ultime_http_cli -command del -key A`






