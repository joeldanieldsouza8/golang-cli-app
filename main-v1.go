package main

import (
	"fmt"
	"time"
)

func main_v1() {
	// Initialize a Todos slice
	todos := Todos{}

	// Add new todos
	todos.Add(Todo{Title: "Learn Go", Completed: false, CreatedAt: time.Now()})
	todos.Add(Todo{Title: "Write Go Code", Completed: false, CreatedAt: time.Now()})
	todos.Add(Todo{Title: "Practice Go", Completed: false, CreatedAt: time.Now()})

	// Create a new Storage instance for saving/loading todos
	storage := NewStorage[Todos]("todos.json")

	// Save the current state of todos to the file
	fmt.Println("Saving todos to file...")
	if err := storage.Save(todos); err != nil {
		fmt.Println("Error saving todos:", err)
	} else {
		fmt.Println("Todos successfully saved to file.")
	}

	// Clear the current todos slice to simulate loading from file
	todos = Todos{}

	// Load todos from the file
	fmt.Println("\nLoading todos from file...")
	if err := storage.Load(&todos); err != nil {
		fmt.Println("Error loading todos:", err)
	} else {
		fmt.Println("Todos successfully loaded from file.")
		todos.Print()
	}

	// Remove the second todo (index 1)
	err := todos.Remove(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("\nAfter removing the todo at index 1:")
		todos.Print()
	}

	// Toggle the completion status of the first todo
	fmt.Println("\nToggling completion status of the first todo:")
	err = todos.ToggleComplete(0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		todos.Print()
	}

	// Edit the title of the first todo
	fmt.Println("\nEditing the title of the first todo:")
	err = todos.EditTitle(0, "Master Go")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		todos.Print()
	}

	// Save the modified todos to the file
	fmt.Println("\nSaving modified todos to file...")
	if err := storage.Save(todos); err != nil {
		fmt.Println("Error saving todos:", err)
	} else {
		fmt.Println("Modified todos successfully saved to file.")
	}
}
