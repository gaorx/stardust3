package sderr

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var (
	// New
	New    = errors.New
	Errorf = errors.Errorf
	Prefix = multierror.Prefix

	// With
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf

	// Cause
	Cause  = errors.Cause
	Unwrap = errors.Unwrap

	// Is/As
	Is = errors.Is
	As = errors.As

	// Multiple
	Append     = multierror.Append
	Flatten    = multierror.Flatten
	ListFormat = multierror.ListFormatFunc
)

type (
	Frame           = errors.Frame
	StackTrace      = errors.StackTrace
	MultipleError   = multierror.Error
	ErrorFormatFunc = multierror.ErrorFormatFunc
	Group           = multierror.Group
)

func Select(err1, err2 error) error {
	if err1 == nil && err2 == nil {
		return nil
	} else if err1 != nil && err2 != nil {
		return Append(err1, err2)
	} else if err1 != nil {
		return err1
	} else {
		return err2
	}
}

func Concurrently(funcs ...func() error) *MultipleError {
	var g Group
	for _, f := range funcs {
		g.Go(f)
	}
	return g.Wait()
}
