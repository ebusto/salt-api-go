package salt

import (
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	Message string
	Status  int
}

func NewError(status int, message string) *Error {
	return &Error{
		Message: message,
		Status:  status,
	}
}

func (e *Error) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "[%d] %s", e.Status, http.StatusText(e.Status))

	if len(e.Message) > 0 {
		fmt.Fprintf(&b, ": %s", e.Message)
	}

	return b.String()
}
