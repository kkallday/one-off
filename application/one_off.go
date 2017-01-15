package application

import (
	"fmt"
	"io"
	"os/exec"
)

type OneOff struct {
	fly               fly
	pipelineConverter pipelineConverter
	writer            io.Writer
}

type fly interface {
	GetPipeline(targetAlias, pipeline string) (string, error)
	SetPathToFly(pathToFly string)
}

type pipelineConverter interface {
	EnvVars(pipeline, job, task string) (string, error)
}

var lookPath = exec.LookPath

func NewOneOff(fly fly, pipelineConverter pipelineConverter, writer io.Writer) OneOff {
	return OneOff{
		fly:               fly,
		pipelineConverter: pipelineConverter,
		writer:            writer,
	}
}

func (o OneOff) Run(inputs OneOffInputs) error {
	if inputs.FlyOverride != "" {
		o.fly.SetPathToFly(inputs.FlyOverride)
	} else {
		pathToFly, err := lookPath("fly")
		if err != nil {
			panic(err)
		}
		o.fly.SetPathToFly(pathToFly)
	}

	pipelineYAML, err := o.fly.GetPipeline(inputs.TargetAlias, inputs.Pipeline)
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %v", err)
	}

	envVars, err := o.pipelineConverter.EnvVars(pipelineYAML, inputs.Job, inputs.Task)
	if err != nil {
		return fmt.Errorf("failed to retrieve pipeline params from pipeline: %v", err)
	}

	script := fmt.Sprintf(`#!/bin/bash -exu
%s

fly -t %s execute --config=REPLACE/ME/PATH/TO/TASK --inputs-from %s/%s`,
		envVars, inputs.TargetAlias, inputs.Pipeline, inputs.Job)

	_, err = o.writer.Write([]byte(script))
	if err != nil {
		return fmt.Errorf("failed to write one-off to stdout: %v", err)
	}

	return nil
}
