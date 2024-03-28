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
	return fmt.Sprintf("%s at %d:%d : %s", c.ErrorType, c.Position.Line, c.Position.Column, c.Message)
}
