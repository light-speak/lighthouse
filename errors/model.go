package errors

type DataloaderError struct {
	Msg string
}

func (e *DataloaderError) Error() string {
	return e.Msg
}

