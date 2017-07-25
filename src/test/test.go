package main

import (
    "web/cipher"
    "fmt"
)

func main () {
    fmt.Println(cipher.Encode("http://localhost:8080", "1212", "5"))
}