package errorx

import (
	"context"
	"fmt"
	"strings"
)

type Error struct {
	HttpCode uint
	BizCode  int64
	Message  string
	Domain   string
	Reason   []string
	Wrapped  error
}

func new() *Error {
	return &Error{
		Reason: make([]string, 0),
	}
}

func resolveOptions(err *Error, options ...Option) *Error {
	for _, option := range options {
		option(err)
	}

	return err
}

// New used to create an errorx.Error with options
func New(options ...Option) *Error {
	return resolveOptions(new(), options...)
}

// NewWithCtx used to create an errorx.Error with context.Context and options.
// Received options will cover the options in context.Context
func NewWithCtx(ctx context.Context, options ...Option) *Error {
	return resolveOptions(new(), append(extractFields(ctx), options...)...)
}

func (e *Error) String() string {
	if isCustomSerializer {
		return serializer(e)
	}

	return e.toString()
}

func (e *Error) toString() string {
	var builder strings.Builder

	if e.Domain != "" {
		builder.WriteString(fmt.Sprintf("[%s] ", e.Domain))
	}

	builder.WriteString(fmt.Sprintf("[HttpCode: %d] ", e.HttpCode))
	builder.WriteString(fmt.Sprintf("[BizCode: %d] ", e.BizCode))
	builder.WriteString(fmt.Sprintf("[Message: %s] ", e.Message))

	if len(e.Reason) > 0 {
		var reasonBuilder strings.Builder

		for _, reason := range e.Reason {
			reasonBuilder.WriteString(fmt.Sprintf("%s;", reason))
		}

		builder.WriteString(fmt.Sprintf("[Reason: %s]", reasonBuilder.String()))
	}

	if e.Wrapped != nil {
		builder.WriteString(fmt.Sprintf("[Wrapped: %s] ", e.Wrapped))
	}

	return builder.String()
}
