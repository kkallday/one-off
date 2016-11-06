package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kkallday/one-off/application"
	"github.com/kkallday/one-off/concourse"
)

func main() {
	pathToFly, err := exec.LookPath("fly")
	if err != nil {
		fail(err)
	}

	fly := concourse.NewFly(pathToFly)
	pipelineConverter := concourse.NewPipelineConverter()
	oneOff := application.NewOneOff(fly, pipelineConverter)

	app := application.New(oneOff)
	err = app.Execute(os.Args[1:])
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
