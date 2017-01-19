package main

import (
	"fmt"
	"os"

	"github.com/kkallday/one-off/application"
	"github.com/kkallday/one-off/concourse"
)

func main() {
	fly := concourse.NewFly()
	pipelineConverter := concourse.NewPipelineConverter()
	oneOff := application.NewOneOff(&fly, pipelineConverter, os.Stdout)

	app := application.New(oneOff)
	err := app.Execute(os.Args[1:])
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
