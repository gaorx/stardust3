package sdhttpua

import (
	"github.com/gaorx/stardust3/sdrand"
)

func ChoiceUA(filters ...Filter) *UA {
	l := FindUA(filters...)
	if len(l) == 0 {
		return nil
	}
	return sdrand.ChoiceOne(l).(*UA)
}

func Choice(filters ...Filter) string {
	ua := ChoiceUA(filters...)
	if ua == nil {
		return ""
	}
	return ua.UA
}
