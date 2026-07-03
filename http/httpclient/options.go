package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Option func(*client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *client) {
		c.timeout = timeout
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(c *client) {
		if cfg == nil {
			return
		}

		httpClient, ok := c.http.(*http.Client)
		if !ok {
			return
		}

		transport, ok := httpClient.Transport.(*http.Transport)
		if !ok || transport == nil {
			defaultTransport, ok := http.DefaultTransport.(*http.Transport)
			if !ok {
				return
			}

			transport = defaultTransport.Clone()
		}

		transport.TLSClientConfig = cfg
		httpClient.Transport = transport
	}
}

type RequestOption func(*Request)

func WithAcceptStatusCode(code int) RequestOption {
	return func(r *Request) {
		if r.acceptStatusCodes == nil {
			r.acceptStatusCodes = make(map[int]bool, 0)
		}

		r.acceptStatusCodes[code] = true
	}
}
