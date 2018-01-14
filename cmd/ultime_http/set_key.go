package main

import (
	"encoding/json"
	"fmt"

	"github.com/Boomerangz/ultime/cache"
	"github.com/valyala/fasthttp"
)

func SetByKey(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	key := ctx.UserValue("key").(string)
	body := ctx.PostBody()
	var parsedBody map[string]interface{}
	json.Unmarshal(body, &parsedBody)
	err := cache.SetByKey(name, key, parsedBody["value"])
	if err == nil {
		ctx.SetStatusCode(fasthttp.StatusCreated)
		return
	} else if err == cache.NotFoundError {
		ctx.SetStatusCode(404)
		fmt.Fprint(ctx, err)
		return
	} else {
		ctx.SetStatusCode(500)
		fmt.Fprint(ctx, err)
		return
	}
}
