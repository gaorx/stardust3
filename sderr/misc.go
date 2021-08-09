package sderr

func ToErr(v interface{}) error {
	switch err := v.(type) {
	case nil:
		return nil
	case error:
		return err
	case string:
		return New(err)
	default:
		return Newf("%v", err)
	}
}
