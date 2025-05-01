package errorxv1

import "errors"

func As(err error) *Error {
	if err == nil {
		return nil
	}

	var e *Error
	ok := errors.As(err, &e)
	if !ok {
		e.HttpCode = 500
		e.BizCode = -1
		e.Message = "unknown error"
	}

	return e
}
