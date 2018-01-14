package main

import (
	"fmt"

	"github.com/Boomerangz/ultime/cache"
	"github.com/valyala/fasthttp"
)

func Remove(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	err := cache.Remove(name)
	if err == cache.NotFoundError {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		fmt.Fprint(ctx, err)
	} else if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprint(ctx, err)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}
}
