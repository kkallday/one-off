package application

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type OneOff struct {
	fly               fly
	pipelineConverter pipelineConverter
}

type fly interface {
	GetPipeline(targetAlias, pipeline string) (string, error)
}

type pipelineConverter interface {
	EnvVars(pipeline, job, task string) (string, error)
}

func NewOneOff(fly fly, pipelineConverter pipelineConverter) OneOff {
	return OneOff{
		fly:               fly,
		pipelineConverter: pipelineConverter,
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

	var dir string
	if inputs.OutputDir != "" {
		dir = inputs.OutputDir
	}

	err = ioutil.WriteFile(filepath.Join(dir, fmt.Sprintf("%s-one-off", inputs.Task)), []byte(script), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write script: %v", err)
	}

	return nil
}
