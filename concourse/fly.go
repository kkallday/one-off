package concourse

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Fly struct {
	pathToFly string
}

func NewFly(pathToFly string) Fly {
	return Fly{
		pathToFly: pathToFly,
	}
}

func (f Fly) GetPipeline(targetAlias, pipeline string) (string, error) {
	cmd := exec.Command(f.pathToFly, "-t", "some-target", "get-pipeline", "--pipeline", "some-pipeline")

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	switch err.(type) {
	case *exec.ExitError:
		return "", fmt.Errorf("%v, stderr: %s", err, stderr.Bytes())
	case error:
		return "", err
	}

	return string(stdout.Bytes()), nil
}
