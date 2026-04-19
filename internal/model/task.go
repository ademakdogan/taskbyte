package model

import "time"

// Status represents a task's current state.
type Status string

const (
StatusTodo       Status = "todo"
StatusInProgress Status = "in_progress"
StatusDone       Status = "done"
StatusCancelled  Status = "cancelled"
)

// NextStatus returns the next status in the cycle:
// Todo -> InProgress -> Done -> Cancelled -> Todo
func (s Status) NextStatus() Status {
switch s {
case StatusTodo:
return StatusInProgress
case StatusInProgress:
return StatusDone
case StatusDone:
return StatusCancelled
case StatusCancelled:
return StatusTodo
default:
return StatusTodo
}
}

// String returns the display string for the status.
func (s Status) String() string {
return string(s)
}

// Symbol returns the checkbox symbol for the status.
func (s Status) Symbol() string {
switch s {
case StatusTodo:
return "[ ]"
case StatusInProgress:
return "[…]"
case StatusDone:
return "[✓]"
case StatusCancelled:
return "[x]"
default:
return "[ ]"
}
}

// Label returns the display label (e.g., "@in_progress").
func (s Status) Label() string {
switch s {
case StatusInProgress:
return "@in_progress"
case StatusDone:
return "@done"
case StatusCancelled:
return "@cancelled"
default:
return ""
}
}

// Task represents a single todo item.
type Task struct {
ID              int
Text            string
Status          Status
Date            string // stored as YYYY-MM-DD
CreatedAt       time.Time
UpdatedAt       time.Time
StatusChangedAt *time.Time
}
