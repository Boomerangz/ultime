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
	router.GET("/cache/:name", HttpLog(Get))
	router.GET("/cache/:name/:key", HttpLog(GetByKey))
	router.POST("/cache/:name", HttpLog(Set))
	router.POST("/cache/:name/:key", HttpLog(SetByKey))
	router.DELETE("/cache/:name", HttpLog(Remove))
	router.GET("/cache-keys/", HttpLog(Keys))
	port := config.GetConfig().Port
	log.Printf("Starting server and %d\n", port)

	listenForShutDown()
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), router.Handler))
}

func HttpLog(inFunc func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		inFunc(ctx)
		log.Printf("%s %s %d %d bytes\n", string(ctx.Method()), string(ctx.Path()), ctx.Response.StatusCode(), len(ctx.Response.Body()))
	}
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
