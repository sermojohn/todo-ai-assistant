package main

import (
    "bytes"
    "io"
    "os"
    "strings"
    "testing"
    "time"
)

func TestTimeFormat(t *testing.T) {
    // Ensure time.Now() formatting doesn't panic and returns something parseable
    now := time.Now().Format(time.RFC3339)
    if _, err := time.Parse(time.RFC3339, now); err != nil {
        t.Fatalf("time format invalid: %v", err)
    }
}

// captureOutput runs f while capturing stdout and returns the captured string.
func captureOutput(f func()) string {
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    // run the function while stdout is redirected
    f()

    // restore stdout and collect output
    w.Close()
    var buf bytes.Buffer
    io.Copy(&buf, r)
    os.Stdout = old
    return buf.String()
}

func TestMainOutput(t *testing.T) {
    tests := []struct{
        name string
        want string
    }{
        {name: "contains hello world", want: "hello world!"},
    }

    for _, tc := range tests {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            out := captureOutput(main)
            if !strings.Contains(out, tc.want) {
                t.Fatalf("output = %q, want substring %q", out, tc.want)
            }

            // also assert the output begins with an RFC3339 timestamp
            fields := strings.Fields(out)
            if len(fields) == 0 {
                t.Fatalf("no output captured")
            }
            if _, err := time.Parse(time.RFC3339, fields[0]); err != nil {
                t.Fatalf("timestamp invalid: %v; output=%q", err, out)
            }
        })
    }
}
