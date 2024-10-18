package main

import (
	"encoding/json"
	"os"
)

type Storage[T any] struct {
	fileName string
}

// This function (must) returns a pointer to a new Storage instance struct
func NewStorage[T any](fileName string) *Storage[T] {
	// Create a new of 'Storage[T]' and initialize it with the given file name
	// The (address-of operator) '&' is used to get the memory address of the struct
	// Then, the memory address of the instance is returned to the caller as a pointer to the struct instance
	// The caller can then use this pointer to access the struct instance
	// In summary, this returns the pointer to the newly created 'Storage[T]' struct instance. By using '&', you're returning a memory reference (pointer) instead of a copy of the new struct instance.
	return &Storage[T]{fileName: fileName} // The '&' operator turns it into a pointer

	// This is the same as:
	// storage := &Storage[T]{fileName: fileName} // 'storage' is a pointer to the new 'Storage[T]' struct instance. Hence, points to the memory address of the struct instance
	// return storage // Return the pointer to the new 'Storage[T]' struct instance
}

func (s *Storage[T]) Save(data T) error {
	fileData, error := json.Marshal(data)

	// Check if there was an error while marshalling the data
	if error != nil {
		return error
	}

	// Write the marshalled data to the file
	return os.WriteFile(s.fileName, fileData, 0644)
}

func (s *Storage[T]) Load(data *T) error {
	// Read the file data
	fileData, error := os.ReadFile(s.fileName)

	// Check if there was an error while reading the file
	if error != nil {
		return error
	}

	// Unmarshal the file data into the 'data' pointer
	return json.Unmarshal(fileData, data)
}
