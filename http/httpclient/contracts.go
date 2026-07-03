package httpclient

import "context"

type Client interface {
	Patch(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error)
	Put(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error)
	Post(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error)
	PostForm(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error)
	Delete(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error)
	Get(ctx context.Context, request *Request, opts ...RequestOption) (rst Result, err error)
}
