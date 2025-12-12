# Todo CLI

Small Go CLI for managing a simple todo list stored as JSON.

Build:

```sh
go build -o todo
```

Run examples:

```sh
# Add a task
./todo add "Buy milk"

# List pending tasks
./todo list

# List all tasks including done ones
./todo list --all

# List tasks by priority from high to low
./todo list --priority 

# Mark task 1 done
./todo done 1

# Remove task 2
./todo rm 2
```

Test:

```sh
go test ./...
```

Notes:

- By default tasks are persisted to `$HOME/.todo.json`. Override with `TODO_FILE=/path/to/file ./todo ...`.
- This implementation uses only the Go standard library and stores tasks as an array of objects in JSON.
