package errorxv1

type Option func(err *Error)

func WithHttpCode(code int) Option {
	return func(err *Error) {
		err.HttpCode = code
	}
}

func WithBizCode(code int64) Option {
	return func(err *Error) {
		err.BizCode = code
	}
}

func WithMessage(message string) Option {
	return func(err *Error) {
		err.Message = message
	}
}

func WithDomain(domain string) Option {
	return func(err *Error) {
		err.Domain = domain
	}
}

func WithReason(reason string) Option {
	return func(err *Error) {
		err.Reason = append(err.Reason, reason)
	}
}

func WithError(e error) Option {
	return func(err *Error) {
		err.Wrapped = e
	}
}
