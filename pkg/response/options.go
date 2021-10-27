package response

import "net/http"

type Options func(opt *Opt)

type Opt struct {
	Code    int
	Headers map[string]string
}

func WithCode(c int) Options {
	return func(opt *Opt) {
		opt.Code = c
	}
}

func WithHeaders(h map[string]string) Options {
	return func(opt *Opt) {
		opt.Headers = h
	}
}

func NewOptions() *Opt {
	return &Opt{
		Code: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":           "application/json;charset=utf8",
			"X-Content-Type-Options": "nosniff",
		},
	}
}
