package newsAgr

type ValidationError struct {
	err    error
	reason string
}

func (v ValidationError) Unwrap() error {
	return v.err
}

func (v ValidationError) Error() string {
	return v.err.Error() + ": " + v.reason
}
