package main

import (
	"context"
	"encoding/json"

	"github.com/lcnascimento/go-kit/httpclient"
	"github.com/lcnascimento/go-kit/log"
)

func main() {
	ctx := context.Background()

	log.SetLevel(log.LevelDebug)
	client := httpclient.New()

	res, err := client.Get(ctx, &httpclient.HTTPRequest{
		Host: "https://api.open-meteo.com",
		Path: "/v1/forecast",
		QueryParams: httpclient.HTTPQueryParams{
			"latitude":  "42.22",
			"longitude": "23.39",
			"current":   "temperature_2m",
		},
	})
	if err != nil {
		log.Fatal(ctx, err)
	}

	var data map[string]any
	if err := json.Unmarshal(res.Response, &data); err != nil {
		log.Fatal(ctx, err)
	}

	log.Info(ctx, "temperature fetched", log.Any("data", data))
}
