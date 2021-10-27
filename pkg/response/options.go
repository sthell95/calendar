package response

import "net/http"

type Options func(opt *Opt)

type Opt struct {
	Code    int
	Message string
	Headers map[string]string
}

func WithCode(c int) Options {
	return func(opt *Opt) {
		opt.Code = c
	}
}

func WithMessage(m string) Options {
	return func(opt *Opt) {
		opt.Message = m
	}
}

func WithHeaders(h map[string]string) Options {
	return func(opt *Opt) {
		opt.Headers = h
	}
}

func (o *Opt) init() {
	o.Code = http.StatusOK
	o.Headers = map[string]string{
		"Content-Type":           "application/json;charset=utf8",
		"X-Content-Type-Options": "nosniff",
	}
}
