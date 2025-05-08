package response

import "fmt"

type InvalidOrderResponseWriter struct {
	expectedState writerState
	actual        writerState
}

func (e *InvalidOrderResponseWriter) Error() string {
	err := fmt.Sprintf("invalid order while writing response: expected = %s\n\tactual = %s",&e.expectedState, &e.actual)
	return err
}
