package sdtaskcron

import (
	"strings"
)

func ActionIs(actionID string) TaskPred {
	return func(t *Task) bool {
		if t == nil {
			return false
		}
		return t.Action == actionID
	}
}

func ActionHasPrefix(prefix string) TaskPred {
	return func(t *Task) bool {
		if t == nil {
			return false
		}
		if prefix == "" {
			return true
		}
		return strings.HasPrefix(t.Action, prefix)
	}
}
