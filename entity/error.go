package entity

import "fmt"

type CustomError struct {
	ErrorType string
	Message   string
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", c.ErrorType, c.Message)
}
