package todo

import (
	"errors"
	"time"
)

// Task represents a todo item.
type Task struct {
	ID       int64     `json:"id"`
	Text     string    `json:"text"`
	Done     bool      `json:"done"`
	Created  time.Time `json:"created"`
	Priority int       `json:"priority,omitempty"` // 3=high,2=med,1=low
}

// AddTask adds a new task to the given file (if path=="" uses default path).
func AddTask(path, text string) (Task, error) {
	if text == "" {
		return Task{}, errors.New("empty task")
	}
	p := filePathOrDefault(path)
	tasks, err := loadTasks(p)
	if err != nil {
		return Task{}, err
	}
	var maxID int64
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	nt := Task{
		ID:      maxID + 1,
		Text:    text,
		Done:    false,
		Created: time.Now(),
	}
	tasks = append(tasks, nt)
	if err := saveTasks(p, tasks); err != nil {
		return Task{}, err
	}
	return nt, nil
}

// ListTasks returns all tasks from storage at path (or default if path=="").
func ListTasks(path string) ([]Task, error) {
	p := filePathOrDefault(path)
	return loadTasks(p)
}

// MarkDone marks a task as done by id.
func MarkDone(path string, id int64) error {
	p := filePathOrDefault(path)
	tasks, err := loadTasks(p)
	if err != nil {
		return err
	}
	var found bool
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}
	if !found {
		return errors.New("task not found")
	}
	return saveTasks(p, tasks)
}

// RemoveTask deletes a task by id.
func RemoveTask(path string, id int64) error {
	p := filePathOrDefault(path)
	tasks, err := loadTasks(p)
	if err != nil {
		return err
	}
	var out []Task
	var found bool
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		out = append(out, t)
	}
	if !found {
		return errors.New("task not found")
	}
	return saveTasks(p, out)
}

// SetPriority sets the numeric priority for a task by id (3=high,2=med,1=low).
func SetPriority(path string, id int64, priority int) error {
	if priority < 1 || priority > 3 {
		return errors.New("invalid priority")
	}
	p := filePathOrDefault(path)
	tasks, err := loadTasks(p)
	if err != nil {
		return err
	}
	var found bool
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Priority = priority
			found = true
			break
		}
	}
	if !found {
		return errors.New("task not found")
	}
	return saveTasks(p, tasks)
}

// filePathOrDefault picks the provided path or falls back to TodoFilePath().
func filePathOrDefault(path string) string {
	if path != "" {
		return path
	}
	return TodoFilePath()
}
