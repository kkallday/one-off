package fakes

type Fly struct {
	GetPipelineCall struct {
		CallCount int
		Receives  struct {
			TargetAlias string
			Pipeline    string
		}
		Returns struct {
			PipelineYAML string
			Error        error
		}
	}
}

func (f *Fly) GetPipeline(targetAlias string, pipeline string) (string, error) {
	f.GetPipelineCall.CallCount++
	f.GetPipelineCall.Receives.TargetAlias = targetAlias
	f.GetPipelineCall.Receives.Pipeline = pipeline
	return f.GetPipelineCall.Returns.PipelineYAML, f.GetPipelineCall.Returns.Error
}
