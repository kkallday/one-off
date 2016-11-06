package fakes

type PipelineConverter struct {
	EnvVarsCall struct {
		CallCount int
		Receives  struct {
			PipelineYAML string
			Job          string
			Task         string
		}
		Returns struct {
			EnvVars string
			Error   error
		}
	}
}

func (p *PipelineConverter) EnvVars(pipelineYAML, job, task string) (string, error) {
	p.EnvVarsCall.CallCount++
	p.EnvVarsCall.Receives.PipelineYAML = pipelineYAML
	p.EnvVarsCall.Receives.Job = job
	p.EnvVarsCall.Receives.Task = task
	return p.EnvVarsCall.Returns.EnvVars, p.EnvVarsCall.Returns.Error
}
