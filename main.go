package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/example/todo/internal/todo"
)

func usage() {
	prog := filepath.Base(os.Args[0])
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s add <task text>       Add a new task\n", prog)
	fmt.Printf("  %s list [--all] [--top N] List tasks (default: pending only)\n", prog)
	fmt.Printf("  %s done <id>            Mark task <id> done\n", prog)
	fmt.Printf("  %s rm <id>              Remove task <id>\n", prog)
	fmt.Printf("  %s prioritize           Ask GenAI to prioritize tasks\n", prog)
	fmt.Printf("  %s reprioritize <id> <high|medium|low|3|2|1>  Manually set task priority\n", prog)
	fmt.Printf("  %s help                 Show this help\n", prog)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "add: missing task text")
			os.Exit(1)
		}
		text := strings.Join(os.Args[2:], " ")
		t, err := todo.AddTask("", text)
		if err != nil {
			fmt.Fprintln(os.Stderr, "add error:", err)
			os.Exit(1)
		}
		fmt.Printf("added %d: %s\n", t.ID, t.Text)

	case "list":
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		all := fs.Bool("all", false, "show all tasks, including done ones")
		top := fs.Int("top", 0, "show top N tasks by priority")
		n := fs.Int("n", 0, "alias for --top")
		prioritySort := fs.Bool("priority", false, "sort tasks by priority high->low")
		_ = fs.Parse(os.Args[2:])
		tasks, err := todo.ListTasks("")
		if err != nil {
			fmt.Fprintln(os.Stderr, "list error:", err)
			os.Exit(1)
		}

		// filter according to --all
		var filtered []todo.Task
		for _, t := range tasks {
			if !*all && t.Done {
				continue
			}
			filtered = append(filtered, t)
		}

		// determine top N (n overrides top if provided)
		nn := *top
		if *n > 0 {
			nn = *n
		}
		// sort either when user requested priority sort or when requesting top N
		if *prioritySort || nn > 0 {
			sort.Slice(filtered, func(i, j int) bool {
				if filtered[i].Priority != filtered[j].Priority {
					return filtered[i].Priority > filtered[j].Priority
				}
				return filtered[i].Created.Before(filtered[j].Created)
			})
			if nn > 0 && nn < len(filtered) {
				filtered = filtered[:nn]
			}
		}

		for _, t := range filtered {
			mark := ' '
			if t.Done {
				mark = 'x'
			}
			pri := ""
			switch t.Priority {
			case 3:
				pri = "(high)"
			case 2:
				pri = "(med)"
			case 1:
				pri = "(low)"
			}
			fmt.Printf("[%c] %d: %s %s\n", mark, t.ID, t.Text, pri)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "done: missing id")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "done: invalid id")
			os.Exit(1)
		}
		if err := todo.MarkDone("", id); err != nil {
			fmt.Fprintln(os.Stderr, "done error:", err)
			os.Exit(1)
		}
		fmt.Printf("marked %d done\n", id)

	case "rm":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "rm: missing id")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "rm: invalid id")
			os.Exit(1)
		}
		if err := todo.RemoveTask("", id); err != nil {
			fmt.Fprintln(os.Stderr, "rm error:", err)
			os.Exit(1)
		}
		fmt.Printf("removed %d\n", id)

	case "prioritize":
		if err := todo.PrioritizeTasks(""); err != nil {
			fmt.Fprintln(os.Stderr, "prioritize error:", err)
			os.Exit(1)
		}
		fmt.Println("prioritization with AI complete")

	case "reprioritize":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "reprioritize: missing id or priority")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reprioritize: invalid id")
			os.Exit(1)
		}
		pstr := strings.ToLower(os.Args[3])
		var pri int
		switch pstr {
		case "high", "3":
			pri = 3
		case "medium", "med", "2":
			pri = 2
		case "low", "1":
			pri = 1
		default:
			// try numeric parse
			v, err := strconv.Atoi(pstr)
			if err != nil || v < 1 || v > 3 {
				fmt.Fprintln(os.Stderr, "reprioritize: invalid priority, use high|medium|low or 3|2|1")
				os.Exit(1)
			}
			pri = v
		}
		if err := todo.SetPriority("", id, pri); err != nil {
			fmt.Fprintln(os.Stderr, "reprioritize error:", err)
			os.Exit(1)
		}
		fmt.Printf("set %d priority to %d\n", id, pri)

	case "help":
		usage()

	default:
		fmt.Fprintln(os.Stderr, "unknown command:", cmd)
		usage()
		os.Exit(1)
	}
}
