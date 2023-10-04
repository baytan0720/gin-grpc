package gin_grpc

type Errors struct {
	Errors []error
}

func (e *Errors) Last() error {
	if length := len(e.Errors); length > 0 {
		return e.Errors[length-1]
	}
	return nil
}

func (e *Errors) Append(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *Errors) Len() int {
	return len(e.Errors)
}

func (e *Errors) Clear() {
	e.Errors = e.Errors[:0]
}
