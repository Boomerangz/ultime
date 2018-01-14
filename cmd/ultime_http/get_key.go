package main

import (
	"encoding/json"
	"fmt"

	"github.com/Boomerangz/ultime/cache"
	"github.com/valyala/fasthttp"
)

func GetByKey(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	key := ctx.UserValue("key").(string)
	value, err := cache.GetByKey(name, key)
	if err == nil && value != nil {
		result := map[string]interface{}{}
		result["value"] = value
		marshalled, err := json.Marshal(result)
		if err == nil {
			fmt.Fprint(ctx, string(marshalled))
			return
		}
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
