package sdstr

import (
	"strings"
)

func Split2s(s, sep string) (string, string) {
	if s == "" {
		return "", ""
	}
	l := strings.SplitN(s, sep, 2)
	switch len(l) {
	case 0:
		return "", ""
	case 1:
		return l[0], ""
	default:
		return l[0], l[1]
	}
}

func Split3s(s, sep string) (string, string, string) {
	if s == "" {
		return "", "", ""
	}
	l := strings.SplitN(s, sep, 3)
	switch len(l) {
	case 0:
		return "", "", ""
	case 1:
		return l[0], "", ""
	case 2:
		return l[0], l[1], ""
	default:
		return l[0], l[1], l[2]
	}
}

func Split4s(s, sep string) (string, string, string, string) {
	if s == "" {
		return "", "", "", ""
	}
	l := strings.SplitN(s, sep, 4)
	switch len(l) {
	case 0:
		return "", "", "", ""
	case 1:
		return l[0], "", "", ""
	case 2:
		return l[0], l[1], "", ""
	case 3:
		return l[0], l[1], l[2], ""
	default:
		return l[0], l[1], l[2], l[3]
	}
}
