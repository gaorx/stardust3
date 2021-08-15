package sdecho

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
}

func H(h func(Context) error) func(echo.Context) error {
	if h == nil {
		panic(sderr.New("nil handler"))
	}
	return func(ec0 echo.Context) error {
		return sderr.WithStack(
			h(Context{ec0}),
		)
	}
}

func E(eh func(error, Context)) func(error, echo.Context) {
	if eh == nil {
		panic(sderr.New("nil error handler"))
	}
	return func(err0 error, ec0 echo.Context) {
		eh(err0, Context{ec0})
	}
}
