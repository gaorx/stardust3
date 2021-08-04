package sderr

import (
	"github.com/pkg/errors"
)

var (
	// New
	New    = errors.New
	Errorf = errors.Errorf

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
)

type (
	Frame      = errors.Frame
	StackTrace = errors.StackTrace
)
