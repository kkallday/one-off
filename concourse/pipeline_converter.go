package concourse

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type PipelineConverter struct{}

func NewPipelineConverter() PipelineConverter {
	return PipelineConverter{}
}

type Pipeline struct {
	Jobs []Job
}

type Job struct {
	Name  string
	Plans []Plan `yaml:"plan"`
}

type Plan struct {
	Task   string
  Params map[string]string
}

func (p PipelineConverter) EnvVars(pipelineYAML, jobName, taskName string) (string, error) {
	var pipeline Pipeline
	err := yaml.Unmarshal([]byte(pipelineYAML), &pipeline)
	if err != nil {
		return "", err
	}

	if len(pipeline.Jobs) == 0 {
		return "", errors.New("pipeline does not contain any jobs - are you sure the pipeline exists?")
	}

	job, err := p.findJob(pipeline, jobName)
	if err != nil {
		return "", err
	}

	task, err := p.findTask(job, taskName, jobName)
	if err != nil {
		return "", err
	}

	var envVars []string
	for k, v := range task.Params {
		envVars = append(envVars, fmt.Sprintf("export %s=%q", k, v))
	}

	sort.Strings(envVars)

	return strings.Join(envVars, "\n"), nil
}

func (PipelineConverter) findJob(pipeline Pipeline, jobName string) (Job, error) {
	for _, job := range pipeline.Jobs {
		if job.Name == jobName {
			return job, nil
		}
	}

	return Job{}, fmt.Errorf("could not find job %q in pipeline", jobName)
}

func (PipelineConverter) findTask(job Job, taskName, jobName string) (Plan, error) {
	for _, plan := range job.Plans {
		if plan.Task == taskName {
			return plan, nil
		}
	}

	return Plan{}, fmt.Errorf("could not find task %q in job %q in pipeline", taskName, jobName)
}
