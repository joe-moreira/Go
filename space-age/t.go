package main

import (
  "fmt"
  "strings"
)

func main() {
  fmt.Println(strings.ContainsAny("Hello World", ",|"))
  fmt.Println(strings.ContainsAny("Hello, World", ",|"))
  fmt.Println(strings.ContainsAny("Hello|World", ",|"))
}
