package exception

type NotFoundErr struct {
	Error string
}

func NewNotFoundErr(err string) NotFoundErr {
	return NotFoundErr{
		Error: err,
	}
}
