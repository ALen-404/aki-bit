package apiclient

import "fmt"

type ApiError struct {
	Code int
	Body []byte
}

func (e ApiError) Error() string {
	return fmt.Sprintf("---\n%d\n%s\n---", e.Code, string(e.Body))
}
