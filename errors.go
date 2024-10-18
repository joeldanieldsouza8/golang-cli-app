package main

import "fmt"

type DateParseError struct {
	DateString string
}

func (e *DateParseError) Error() string {
	return fmt.Sprintf("Invalid date format: '%s'. Expected format: YYYY-MM-DD", e.DateString)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type IndexOutOfRangeError struct {
	Index int
}

func (e *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("Index %d is out of range", e.Index)
}

type PriorityParseError struct {
	Priority string
}

func (e *PriorityParseError) Error() string {
	return fmt.Sprintf("Invalid priority level: '%s'. Expected high, medium, or low", e.Priority)
}
