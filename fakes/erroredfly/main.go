package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprint(os.Stderr, "some error message")
	os.Exit(1)
}
