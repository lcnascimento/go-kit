package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/lcnascimento/go-kit/http/httpclient"
	"github.com/lcnascimento/go-kit/o11y/log"
)

var (
	pkg    = "github.com/lcnascimento/go-kit/httpclient/example"
	logger = log.MustNewLogger(pkg)
)

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
		logger.Critical(ctx, err)
		return
	}

	var body map[string]any
	if err := json.Unmarshal(res.Response, &body); err != nil {
		logger.Critical(ctx, err)
		return
	}

	logger.Info(ctx, "temperature fetched", log.Any("data", body))
}
