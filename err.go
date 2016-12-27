package finder

import (
	"bytes"
	"fmt"
)

// MultiErr combines multiple errors.
type MultiErr []error

func (e MultiErr) Error() string {
	buf := &bytes.Buffer{}
	for i, err := range e {
		fmt.Fprintf(buf, "%d. %s\n", i+1, err)
	}
	return buf.String()
}
