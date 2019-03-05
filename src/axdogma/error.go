package axdogma

type wrappedError struct {
	err error
}

func unwrap(err *error) {
	v := recover()

	if rp, ok := v.(wrappedError); ok {
		*err = rp.err
	} else if v != nil {
		panic(v)
	}
}

func wrapAndPanic(err error) {
	panic(wrappedError{err})
}
