package main

import (
	"fmt"
	"github.com/joeldanieldsouza8/golang-cli-app/internal/command"
	"github.com/joeldanieldsouza8/golang-cli-app/internal/storage"
	"github.com/joeldanieldsouza8/golang-cli-app/internal/todo"
	"os"
)

func main() {
	// Parse the command-line flags
	flags := command.ParseFlags()

	// Initialize a Todos slice
	todos := todo.Todos{}

	// Create a new storage instance for saving/loading todos
	newStorage := storage.NewStorage[todo.Todos]("todos.json")
	// Load todos from the file if it exists or create a new Todos slice if it doesn't exist before executing the commands
	if err := newStorage.Load(&todos); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File todos.json not found, initializing an empty todo list.")
		} else {
			fmt.Println("Error loading todos:", err)
			return // Exit if there's any other error (e.g., file permission issue)
		}
	} else {
		fmt.Println("Todos successfully loaded from file.")
	}

	// Execute the command based on the parsed flags
	if err := flags.Execute(&todos); err != nil {
		fmt.Println("Error executing command:", err)
	}

	// Save the current state of todos to the file
	if err := newStorage.Save(todos); err != nil {
		fmt.Println("Error saving todos:", err)
	} else {
		fmt.Println("Todos successfully saved to file.")
	}

	// Print the updated list of todos to verify the addition
	fmt.Println("Updated list of todos:")
	todos.Print()
}

// Testing for debugging purposes
//func main() {
//	// Override the command-line arguments for testing
//	os.Args = []string{"cmd", "-add", "Read Golang documentation:2024-10-17:high"}
//
//	// Initialize a Todos slice
//	todos := Todos{
//		{Title: "Learn Golang", CreatedAt: time.Now(), Priority: High},
//		{Title: "Build a Todo App", CreatedAt: time.Now(), Priority: Medium},
//		{Title: "Master Go Generics", CreatedAt: time.Now(), Priority: Low},
//		{Title: "Contribute to Open Source", CreatedAt: time.Now(), Priority: High},
//		{Title: "Write a Go Article", CreatedAt: time.Now(), Priority: Medium},
//		{Title: "Learn Go Concurrency", CreatedAt: time.Now(), Priority: Low},
//	}
//
//	// Initialize a new storage instance
//	storage := NewStorage[Todos]("debugging_todos.json")
//
//	// Load todos from the file if it exists
//	if err := storage.Load(&todos); err != nil {
//		fmt.Println("Error loading todos from file:", err)
//	} else {
//		fmt.Println("Todos successfully loaded from file.")
//	}
//
//	// Parse the command-line flags
//	flags := ParseFlags()
//
//	// Execute the command based on the parsed flags
//	err := flags.Execute(&todos)
//	if err != nil {
//		fmt.Println("Error executing command:", err)
//	}
//
//	// Save the current state of todos to the file
//	if err := storage.Save(todos); err != nil {
//		fmt.Println("Error saving todos to file:", err)
//	} else {
//		fmt.Println("Todos successfully saved to file.")
//	}
//}
