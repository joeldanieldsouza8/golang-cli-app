package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CommandFlags stores the parsed command-line flags
type CommandFlags struct {
	Add            string
	Remove         int
	ToggleComplete int
	EditTitle      string
	List           bool
	Help           bool
	Completed      bool
	Sort           string
}

// Since we aren't modifying the 'CommandFlags' struct instance, we don't need to pass it as a pointer
func ParseFlags() CommandFlags {
	// Create a new CommandFlags instance
	// flags := &CommandFlags{} // 'flags' is a pointer to the new 'CommandFlags' struct instance memory address

	// Create a new CommandFlags instance
	var flags CommandFlags

	flag.StringVar(&flags.Add, "add", "", "Add a new todo") // '&flags.Add' gets the memory address of the 'Add' field in the 'CommandFlags' struct instance. This is needed because the StringVar function will store the parsed flag value at this memory address
	flag.IntVar(&flags.Remove, "remove", -1, "Remove a todo by index")
	flag.IntVar(&flags.ToggleComplete, "toggle", -1, "Toggle the completion status of a todo by index")
	flag.StringVar(&flags.EditTitle, "edit", "", "Edit the title of a todo by index")
	flag.BoolVar(&flags.List, "list", false, "List all todos")
	flag.BoolVar(&flags.Help, "help", false, "Display help information")
	flag.BoolVar(&flags.Completed, "completed", false, "List the number of all completed todos")
	// flag.BoolVar(&flags.Sort, "sort", false, "Sort the todos by priority and/or due date")
	flag.StringVar(&flags.Sort, "sort", "", "Sort the todos by priority in asc or desc order")

	// Parse the command-line flags
	flag.Parse()

	// Check if the help flag is provided
	if flags.Help {
		PrintHelp()
		os.Exit(0) // Exit the program with a successful status code
	}

	return flags
}

// type Split struct {
// 	Part1 interface{}
// 	Part2 string
// }

// func splitString(s string) (Split, error) {
// 	// Check if the string contains a colon
// 	if !strings.Contains(s, ":") {
// 		return Split{}, fmt.Errorf("Invalid edit format. Expected <int | string>:<string>")
// 	}

// 	// Split the string into 2 parts
// 	parts := strings.SplitN(s, ":", 2)

// 	// Check if there are 2 parts
// 	if len(parts) != 2 {
// 		return Split{}, fmt.Errorf("Invalid edit format. Expected <int | string>:<string>")
// 	}

// 	// Check if the first part can be parsed as an integer
// 	index, err := strconv.Atoi(parts[0])
// 	if err == nil {
// 		return Split{Part1: index, Part2: parts[1]}, nil
// 	}

// 	// Otherwise, if the first part can't be parsed as an integer then return as (string, string)
// 	return Split{Part1: parts[0], Part2: parts[1]}, nil
// }

func (flags *CommandFlags) Execute(todos *Todos) error {
	switch {
	case flags.List:
		return flags.handleList(todos)

	case flags.Add != "":
		return flags.handleAdd(todos)

	case flags.EditTitle != "":
		return flags.handleEdit(todos)

	case flags.Remove != -1:
		return flags.handleRemove(todos)

	case flags.ToggleComplete != -1:
		return flags.handleToggle(todos)

	case flags.Completed:
		return flags.handleCompleted(todos)

	default:
		return &ValidationError{Message: "No command provided. Use -help to see the list of available commands"}
	}
}

func (flags *CommandFlags) handleList(todos *Todos) error {
	todos.List()
	return nil
}

func (flags *CommandFlags) handleAdd(todos *Todos) error {
	// Check if the string contains a colon
	if !strings.Contains(flags.Add, ":") {
		return &ValidationError{Message: "Invalid add format. Expected <string>:<YYYY-MM-DD>:<high | medium | low>"}
	}

	// Split the string into 3 parts
	parts := strings.SplitN(flags.Add, ":", 3)

	// Check if there are 3 parts
	if len(parts) != 3 {
		return &ValidationError{Message: "Invalid add format. Expected <string>:<YYYY-MM-DD>:<high | medium | low>"}
	}

	// Check if the second part can be parsed as the valid date format
	dueDate, err := time.Parse("2006-01-02", parts[1])
	if err != nil {
		return &DateParseError{DateString: parts[1]}
	}

	// Check if the third part is a valid priority level
	priority, err := parsePriority(parts[2])
	if err != nil {
		return &PriorityParseError{Priority: parts[2]}
	}

	todos.Add(Todo{Title: parts[0], DueDate: &dueDate, Priority: priority})

	return nil
}

func (flags *CommandFlags) handleEdit(todos *Todos) error {
	// Check if the string contains a colon
	if !strings.Contains(flags.EditTitle, ":") {
		return &ValidationError{Message: "Invalid edit format. Expected <int>:<string>"}
	}

	// Split the string into 2 parts
	parts := strings.SplitN(flags.EditTitle, ":", 2)

	// Check if there are 2 parts
	if len(parts) != 2 {
		return &ValidationError{Message: "Invalid edit format. Expected <int>:<string>"}
	}

	// Check if the first part can be parsed as an integer
	index, err := strconv.Atoi(parts[0])
	if err != nil {
		return &ValidationError{Message: "Invalid index. Expected an integer"}
	}

	// This is using the scoped approach to handle the error as the variable 'err' is already declared above
	if err := todos.EditTitle(index, parts[1]); err != nil {
		return err
	}

	return nil
}

func (flags *CommandFlags) handleRemove(todos *Todos) error {
	if err := todos.Remove(flags.Remove); err != nil {
		return err
	}
	return nil
}

func (flags *CommandFlags) handleToggle(todos *Todos) error {
	if err := todos.ToggleComplete(flags.ToggleComplete); err != nil {
		return err
	}
	return nil
}

func (flags *CommandFlags) handleCompleted(todos *Todos) error {
	todos.Completed()
	return nil
}

func (flags *CommandFlags) handleSort(todos *Todos) error {
	if flags.Sort != "" {
		switch strings.ToLower(flags.Sort) {
		case "asc":
			todos.SortByPriorityAsc()

		case "desc":
			todos.SortByPriorityDesc()

		default:
			return &ValidationError{Message: "Invalid sort order. Expected asc or desc"}
		}
	}

	return nil
}

func parsePriority(s string) (TPriority, error) {
	switch strings.ToLower(s) {
	case "high":
		return High, nil
	case "medium":
		return Medium, nil
	case "low":
		return Low, nil
	default:
		return "", fmt.Errorf("invalid priority level: %s. Expected high, medium, or low", s)
	}
}

// Displays a list of all available commands and usage
func PrintHelp() {
	// Define a consistent alignment width for commands
	const cmdWidth = -30 // Adjust the width based on your longest command

	fmt.Println("Usage:")
	fmt.Printf("  %-*s : %s\n", cmdWidth, "-add \"Learn Go:2021-12-31\"", "Add a new todo with a due date")
	fmt.Printf("  %-*s : %s\n", cmdWidth, "-remove 1", "Remove a todo by index")
	fmt.Printf("  %-*s : %s\n", cmdWidth, "-toggle 1", "Toggle the completion status of a todo by index")
	fmt.Printf("  %-*s : %s\n", cmdWidth, "-edit 1:New Title", "Edit the title of a todo by index")
	fmt.Printf("  %-*s : %s\n", cmdWidth, "-list", "List all todos")
	fmt.Printf("  %-*s : %s\n", cmdWidth, "-help", "Display help information")
}
