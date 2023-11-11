package main

import (
	"net/http"
	"rl/pkg/config"
	"rl/pkg/middleware"
	"rl/pkg/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	svr := server.New()
	rlMiddleware := middleware.CreateRateLimiterMiddleware(cfg, "localhost:6379")
	err = http.ListenAndServe(":8080", rlMiddleware(svr))
	if err != nil {
		panic(err)
	}
}
