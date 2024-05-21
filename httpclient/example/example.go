package main

import (
	"context"
	"net/http"
	"os"

	"github.com/lcnascimento/go-kit/httpclient"
	"github.com/lcnascimento/go-kit/log"
)

func main() {
	defer os.Exit(0)

	ctx := context.Background()

	log.SetLevel(log.LevelDebug)
	client := httpclient.New()

	req := &httpclient.Request{
		Host: "https://api.open-meteo.com",
		Path: "/v1/forecast",
		QueryParams: httpclient.QueryParams{
			"latitude":  "42.22",
			"longitude": "23.39",
			"current":   "temperature_2m",
		},
	}

	res, err := client.Get(ctx, req, httpclient.WithAcceptStatusCode(http.StatusOK))
	if err != nil {
		return
	}

	var body map[string]any
	if err := res.Body.Cast(ctx, &body); err != nil {
		return
	}

	log.Info(ctx, "temperature fetched", log.Any("data", body))
}
