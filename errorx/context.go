package errorx

import "context"

type errorKey struct{}

func addField(ctx context.Context, option Option) context.Context {
	options := extractFields(ctx)
	options = append(options, option)
	return context.WithValue(ctx, errorKey{}, options)
}

func extractFields(ctx context.Context) []Option {
	value := ctx.Value(errorKey{})
	result := make([]Option, 0)

	if value != nil {
		options, ok := value.([]Option)
		if ok {
			result = append(result, options...)
		}
	}

	return result
}

func AddErrorHttpCode(ctx context.Context, code int) context.Context {
	return addField(ctx, WithHttpCode(code))
}

func AddBizCode(ctx context.Context, code int64) context.Context {
	return addField(ctx, WithBizCode(code))
}

func AddErrorMessage(ctx context.Context, message string) context.Context {
	return addField(ctx, WithMessage(message))
}

func AddErrorDomain(ctx context.Context, domain string) context.Context {
	return addField(ctx, WithDomain(domain))
}

func AddErrorReason(ctx context.Context, reason string) context.Context {
	return addField(ctx, WithReason(reason))
}

func AddWrappedError(ctx context.Context, err error) context.Context {
	return addField(ctx, WithError(err))
}
