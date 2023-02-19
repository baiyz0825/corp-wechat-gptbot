package xerror

import (
	"fmt"
)

type GPTError struct {
	ErrorCode uint
	Message   string
}

func (e *GPTError) Error() string {
	return fmt.Sprintf("Gpt component error <code,msg> : <%d,%s>", e.ErrorCode, e.Message)
}

func NewGPTError(code uint, msg string) GPTError {
	return GPTError{
		ErrorCode: code,
		Message:   msg,
	}
}
