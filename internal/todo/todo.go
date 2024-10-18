package todo

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aquasecurity/table"
)

type TPriority string

const (
	High   TPriority = "high"
	Medium TPriority = "medium"
	Low    TPriority = "low"
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time // This field is a pointer because it can be null
	DueDate     *time.Time
	Priority    TPriority
}

type Todos []Todo

// This is private function because it starts with a lowercase letter
func validateIndex(index int, todos Todos) error {
	if index < 0 || index >= len(todos) {
		return &IndexOutOfRangeError{Index: index}
	}

	return nil
}

func getPriorityOrder(priority TPriority) int {
	switch priority {
	case High:
		return 3
	case Medium:
		return 2
	case Low:
		return 1
	default:
		return 4
	}
}

func (todos *Todos) Add(todo Todo) {
	*todos = append(*todos, todo)
}

func (todos *Todos) Remove(index int) error {
	// NOTE: In the second argument, we pass the value of the pointer, not the pointer itself. This is because 'todos' is already a pointer, so by using the dereference operator (*) we get the value of the pointer, which is the slice of todos
	err := validateIndex(index, *todos)
	if err != nil {
		return err
	}

	// Use the dereference assignment operator (*) to get the value of the pointer
	// '(*todos)[:index]' gives you the part of the slice up to but not including the item at position index
	// '(*todos)[index+1:]' gives you the part of the slice starting after the item at position index, going all the way to the end
	// For example, if '*todos' is [Todo1, Todo2, Todo3, Todo4] and index = 1, then '(*todos)[:index]' will give [Todo1] and '(*todos)[index+1:]' will give [Todo3, Todo4]
	*todos = append((*todos)[:index], (*todos)[index+1:]...) // The ... is used to "unpack" the slice, so that each element in '(*todos)[index+1:]' is passed individually to append()

	return nil
}

func (todos *Todos) ToggleComplete(index int) error {
	err := validateIndex(index, *todos)
	if err != nil {
		return err
	}

	todo := &(*todos)[index] // Get the pointer to the todo at the specified index
	todo.Completed = !todo.Completed

	if todo.Completed {
		now := time.Now()
		todo.CompletedAt = &now
	} else {
		todo.CompletedAt = nil
	}

	return nil
}

func (todos *Todos) EditTitle(index int, title string) error {
	err := validateIndex(index, *todos)
	if err != nil {
		return err
	}

	// We wrap the dereference operator in parentheses to avoid a syntax error. But usually, if we weren't using the dereference operator, we wouldn't need the parentheses, like so: todos[index].Title = title
	(*todos)[index].Title = title

	return nil
}

func (todos *Todos) SortByPriorityAsc() {
	sort.SliceStable(*todos, func(i, j int) bool {
		return getPriorityOrder((*todos)[i].Priority) < getPriorityOrder((*todos)[j].Priority)
	})
}

func (todos *Todos) SortByPriorityDesc() {
	sort.SliceStable(*todos, func(i, j int) bool {
		return getPriorityOrder((*todos)[i].Priority) > getPriorityOrder((*todos)[j].Priority)
	})
}

func (todos *Todos) List() bool {
	// Check if there are any todos
	if len(*todos) == 0 {
		fmt.Println("You have no todos!")
		return false
	}

	todos.Print()
	return true
}

func (todos *Todos) Completed() {
	completed := 0

	for _, todo := range *todos {
		if todo.Completed {
			completed++
		}
	}

	fmt.Printf("You have %d completed todos\n", completed)
}

// Print the todos in a table format
func (todos *Todos) Print() {
	t := table.New(os.Stdout)

	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetHeaders("#", "Title", "Completed", "Created At", "Completed At", "Due Date", "Priority")
	t.SetHeaderStyle(table.StyleBold)
	t.SetLineStyle(table.StyleBlue)

	for i, todo := range *todos {
		completedAt := "❌"
		if todo.CompletedAt != nil {
			completedAt = todo.CompletedAt.Format(time.RFC1123)
		}

		// Check if the todo has a due date
		dueDate := "No due date"
		if todo.DueDate != nil {
			dueDate = todo.DueDate.Format("2006-01-02")
		}

		t.AddRow(fmt.Sprintf("%d", i+1), todo.Title, fmt.Sprintf("%t", todo.Completed), todo.CreatedAt.Format(time.RFC1123), completedAt, dueDate, string(todo.Priority))
	}

	t.Render()
}
