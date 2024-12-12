package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/lcnascimento/go-kit/httpclient"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var logger *log.Logger

func init() {
	logger = log.NewLogger("github.com/lcnascimento/go-kit/httpclient/example")
}

func main() {
	defer os.Exit(0)

	ctx := context.Background()

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
	if err := json.Unmarshal(res.Body, &body); err != nil {
		logger.Fatal(ctx, err)
		return
	}

	logger.Info(ctx, "temperature fetched", log.Any("data", body))
}
