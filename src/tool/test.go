package main

import (
    "path/filepath"
    "fmt"
)

func main() {
    fmt.Println(filepath.FromSlash("/tmp/shootman.ini"))
}