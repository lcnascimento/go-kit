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

	res, err := client.Get(ctx, &httpclient.Request{
		Host: "https://api.open-meteo.com",
		Path: "/v1/forecast",
		QueryParams: httpclient.QueryParams{
			"latitude":  "42.22",
			"longitude": "23.39",
			"current":   "temperature_2m",
		},
	})
	if err != nil {
		log.Fatal(ctx, err)
	}

	var body map[string]any
	if err := json.Unmarshal(res.Body, &body); err != nil {
		log.Fatal(ctx, err)
	}

	log.Info(ctx, "temperature fetched", log.Any("data", body))
}
