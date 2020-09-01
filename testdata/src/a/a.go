package a

type A error // want "implements error"

var _ error = (A)(nil)

type B struct{}

type C struct{} // want "implements error"

var _ error = (*C)(nil)

func (C) Error() string {
	return ""
}

type D struct{} // want "implements error"

var _ error = (*D)(nil)

func (*D) Error() string {
	return ""
}

type E struct { // want "implements error"
	error
}

var _ error = (*E)(nil)

func f() {
	type E error
	type F struct {
	}
	type G struct {
		error
	}
}
