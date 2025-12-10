package main

import (
    "fmt"
    "time"
)

func main() {
    now := time.Now().Format(time.RFC3339)
    fmt.Printf("%s hello world!\n", now)
}
