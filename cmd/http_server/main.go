package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Boomerangz/ultime/cache"
	"github.com/Boomerangz/ultime/cmd/config"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	cache.Init(config.GetConfig())
	router := fasthttprouter.New()
	router.GET("/cache/:name", Get)
	router.GET("/cache/:name/:key", GetByKey)
	router.POST("/cache/:name", Set)
	router.POST("/cache/:name/:key", SetByKey)
	router.DELETE("/cache/:name", Remove)
	router.GET("/cache-keys/", Keys)
	port := config.GetConfig().Port
	log.Printf("Starting server and %d\n", port)

	listenForShutDown()
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), router.Handler))
}

func listenForShutDown() {

	go func() {
		log.Println("Listening signals...")
		c := make(chan os.Signal, 1)

		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Started saving to disk")
		err := cache.CacheInstance.SaveToDisk()
		if err != nil {
			log.Printf("Saving to disk finished with error: %v\n", err)
		} else {
			log.Printf("Saving to disk finished successfully")
		}
		os.Exit(0)
	}()
}
