package errorsx

import (
	"fmt"
	"strings"
)

type AggregatedError interface {
	error
	Errors() []error
}

func Append(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	}

	// you can call errs = Append(errs, nil), which results in a []error{nil}, not an empty list
	// validate that as well

	hasError := false
	for _, err := range errs {
		if err != nil {
			hasError = true
			break
		}
	}
	if !hasError {
		return err
	}

	if err == nil {
		return &aggregatedError{
			errs: errs,
		}
	}

	if aggregated, ok := err.(AggregatedError); ok { //nolint:errorlint // intentionally checking immediate type to avoid merging into a deeply nested aggregate
		return &aggregatedError{
			errs: append(aggregated.Errors(), errs...),
		}
	}

	return &aggregatedError{
		errs: append([]error{err}, errs...),
	}
}

func Errors(err error) []error {
	if aggregated, ok := err.(AggregatedError); ok { //nolint:errorlint // intentionally checking immediate type
		return aggregated.Errors()
	}
	return []error{err}
}

var (
	_ AggregatedError = (*aggregatedError)(nil)
)

type aggregatedError struct {
	errs []error
}

func (a *aggregatedError) Error() string {
	errs := make([]string, len(a.errs))
	for i, err := range a.errs {
		errs[i] = err.Error()
	}
	return strings.Join(errs, "; ")
}

func (a *aggregatedError) Unwrap() []error {
	return a.errs
}

func (a *aggregatedError) Errors() []error {
	return a.errs
}

func (a *aggregatedError) Format(f fmt.State, verb rune) {
	if verb == 'v' && f.Flag('+') {
		for i, err := range a.errs {
			if i > 0 {
				fmt.Fprint(f, "\n")
			}
			fmt.Fprintf(f, "%+v", err)
		}
		return
	}
	fmt.Fprint(f, a.Error())
}
