package responsor

import (
	"github.com/TremblingV5/box/errorx"
	"net/http"
)

type GinHertzContext interface {
	JSON(code int, obj any)
}

type GinHertzOption func(r *Response)

func WithBizCode(code int64) GinHertzOption {
	return func(r *Response) {
		r.BizCode = code
	}
}

func WithMessage(message string) GinHertzOption {
	return func(r *Response) {
		r.Message = message
	}
}

func AddReason(reason ...string) GinHertzOption {
	return func(r *Response) {
		r.Reasons = append(r.Reasons, reason...)
	}
}

func WithData(data any) GinHertzOption {
	return func(r *Response) {
		r.Data = data
	}
}

func WithDomain(domain string) GinHertzOption {
	return func(r *Response) {
		r.Domain = domain
	}
}

func HttpResp(c GinHertzContext, code int, options ...GinHertzOption) {
	r := &Response{}
	for _, option := range options {
		option(r)
	}
	c.JSON(code, r)
}

func HttpSuccess(c GinHertzContext) {
	HttpResp(
		c,
		http.StatusOK,
		WithMessage(Success),
	)
}

func HttpSuccessWithData(c GinHertzContext, data any) {
	HttpResp(
		c,
		http.StatusOK,
		WithData(data),
		WithMessage(Success),
		WithBizCode(0),
	)
}

func HttpWithErr(c GinHertzContext, err error) {
	e := errorx.As(err)
	HttpResp(
		c,
		e.HttpCode,
		WithMessage(e.Message),
		AddReason(e.Reason...),
		WithDomain(e.Domain),
		WithBizCode(e.BizCode),
	)
}
