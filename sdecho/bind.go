package sdecho

import (
	"io/ioutil"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdjson"
)

func (ec Context) BindJsonValue() (interface{}, error) {
	data, err := ioutil.ReadAll(ec.Request().Body)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	r, err := sdjson.UnmarshalValue(data)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}
