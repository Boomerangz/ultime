package main

import (
	"encoding/json"
	"fmt"

	"github.com/Boomerangz/ultime/cache"
	"github.com/valyala/fasthttp"
)

func Keys(ctx *fasthttp.RequestCtx) {
	keyPattern := string(ctx.QueryArgs().Peek("pattern"))
	keys, err := cache.Keys(keyPattern)
	if err == nil {
		result := map[string]interface{}{}
		result["keys"] = keys
		marshalled, err := json.Marshal(result)
		if err == nil {
			fmt.Fprint(ctx, string(marshalled))
			return
		}
	} else {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, err)
	}
}
