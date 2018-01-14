package main

import (
	"encoding/json"

	"github.com/Boomerangz/ultime/cache"
	"github.com/valyala/fasthttp"
)

func Set(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	body := ctx.PostBody()
	var parsedBody map[string]interface{}
	json.Unmarshal(body, &parsedBody)
	expiresInterface := parsedBody["expires"]
	var expires int
	if expiresInterface != nil {
		expires = int(expiresInterface.(float64))
	}
	cache.Set(name, parsedBody["value"], expires)
	ctx.SetStatusCode(fasthttp.StatusCreated)
}
