package errors

import "fmt"

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("[config error]: %s", e.Message)
}
