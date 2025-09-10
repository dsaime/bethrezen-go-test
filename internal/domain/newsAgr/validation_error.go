package newsAgr

// ValidationError представляет собой ошибку валидации поля, с указанием причины
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
