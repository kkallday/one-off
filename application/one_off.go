package application

import (
	"fmt"
	"io"
)

type OneOff struct {
	fly               fly
	pipelineConverter pipelineConverter
	writer            io.Writer
}

type fly interface {
	GetPipeline(targetAlias, pipeline string) (string, error)
}

type pipelineConverter interface {
	EnvVars(pipeline, job, task string) (string, error)
}

func NewOneOff(fly fly, pipelineConverter pipelineConverter, writer io.Writer) OneOff {
	return OneOff{
		fly:               fly,
		pipelineConverter: pipelineConverter,
		writer:            writer,
	}
}

func (o OneOff) Run(inputs OneOffInputs) error {
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
