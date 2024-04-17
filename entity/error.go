package entity

import (
	"fmt"
)

type CustomError struct {
	ErrorType string
	Message   string
	Position  Position
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%s: %s at %d line", c.ErrorType, c.Message, c.Position.Line)
}
