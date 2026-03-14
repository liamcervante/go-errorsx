package errorsx

import (
	"fmt"
)

type WrappedError interface {
	error
	Unwrap() error
}

func Wrap(err error, format string) error {
	if err == nil {
		return nil
	}
	return New(ErrorCode(err), err, format)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

func Unwrap(err error) error {
	if wrapped, ok := err.(WrappedError); ok { //nolint:errorlint // intentionally checking immediate type
		return wrapped.Unwrap()
	}
	return nil
}

var (
	_ WrappedError = (*wrappedError)(nil)
)

type wrappedError struct {
	current error
	wrapped error
}

func (w *wrappedError) Error() string {
	if w.wrapped == nil {
		return w.current.Error()
	}
	return fmt.Sprintf("%v (%v)", w.current, w.wrapped)
}

func (w *wrappedError) Unwrap() error {
	return w.wrapped
}

func (w *wrappedError) Format(f fmt.State, verb rune) {
	if verb == 'v' && f.Flag('+') {
		fmt.Fprintf(f, "%+v", w.current)
		if w.wrapped != nil {
			fmt.Fprintf(f, "\n  - %+v", w.wrapped)
		}
		return
	}
	fmt.Fprint(f, w.Error())
}
