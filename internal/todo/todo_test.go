package todo

import (
	"path/filepath"
	"testing"
)

func TestAddListMarkRemove(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")

	// Add two tasks
	t1, err := AddTask(path, "first task")
	if err != nil {
		t.Fatalf("AddTask 1 failed: %v", err)
	}
	t2, err := AddTask(path, "second task")
	if err != nil {
		t.Fatalf("AddTask 2 failed: %v", err)
	}

	tasks, err := ListTasks(path)
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// Mark first as done
	if err := MarkDone(path, t1.ID); err != nil {
		t.Fatalf("MarkDone failed: %v", err)
	}
	tasks, _ = ListTasks(path)
	var doneCount int
	for _, tt := range tasks {
		if tt.Done {
			doneCount++
		}
	}
	if doneCount != 1 {
		t.Fatalf("expected 1 done task, got %d", doneCount)
	}

	// Remove second
	if err := RemoveTask(path, t2.ID); err != nil {
		t.Fatalf("RemoveTask failed: %v", err)
	}
	tasks, _ = ListTasks(path)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task after remove, got %d", len(tasks))
	}
}

func TestAddTask_TableDriven(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{"empty text", "", true},
		{"normal text", "something todo", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AddTask(path, tt.text)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestListTasks_TableDriven(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	_, _ = AddTask(path, "first")
	_, _ = AddTask(path, "second")

	tasks, err := ListTasks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("want 2 tasks, got %d", len(tasks))
	}
}

func TestMarkDone_TableDriven(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	t1, _ := AddTask(path, "do task")
	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{"existing", t1.ID, false},
		{"not-found", 9999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MarkDone(path, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestRemoveTask_TableDriven(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	t1, _ := AddTask(path, "remove me")
	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{"existing", t1.ID, false},
		{"not-found", 555, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RemoveTask(path, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestSetPriority_TableDriven(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	t1, _ := AddTask(path, "prio me")
	tests := []struct {
		name     string
		id       int64
		priority int
		wantErr  bool
	}{
		{"valid-high", t1.ID, 3, false},
		{"valid-low", t1.ID, 1, false},
		{"invalid-id", 999, 2, true},
		{"invalid-priority", t1.ID, 99, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetPriority(path, tt.id, tt.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %s failed, expected error=%v, got %v", tt.name, tt.wantErr, err)
			}
		})
	}
}
