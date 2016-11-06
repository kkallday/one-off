package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Printf("%s", strings.Join(os.Args[1:], " "))
}
