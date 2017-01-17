package concourse_test

import (
	"github.com/kkallday/one-off/concourse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	oldPipelineYAML = `
groups:
- name: some-pipeline-name
  jobs:
  - some-job
resources:
- name: some-resource
  type: git
jobs:
- name: some-job
  plan:
  - aggregate:
    - get: some-resource
  - task: some-task
    file: /path/to/task.yml
    config:
      params:
        random-non-param:
        - something-non-param
        VAR1OLD: value1old
        VAR2OLD: value2old
        VAR3OLD: value3old`
	modernPipelineYAML = `---
groups:
- name: some-pipeline-name
  jobs:
  - some-job
resources:
- name: some-resource
  type: git
jobs:
- name: some-job
  plan:
  - aggregate:
    - get: some-resource
  - task: some-task
    file: /path/to/task.yml
    params:
      random-non-param:
      - something-non-param
      VAR1: value1
      VAR2: value2
      VAR3: value3`
)

var _ = Describe("pipeline converter", func() {
	var pc concourse.PipelineConverter

	BeforeEach(func() {
		pc = concourse.NewPipelineConverter()
	})

	Describe("EnvVars", func() {
		It("returns env vars", func() {
			envVars, err := pc.EnvVars(modernPipelineYAML, "some-job", "some-task")
			Expect(err).NotTo(HaveOccurred())
			Expect(envVars).To(Equal(`export VAR1="value1"
export VAR2="value2"
export VAR3="value3"`))
		})

		Context("when given pipeline YAML from older fly versions", func() {
			It("returns env vars", func() {
				envVars, err := pc.EnvVars(oldPipelineYAML, "some-job", "some-task")
				Expect(err).NotTo(HaveOccurred())
				Expect(envVars).To(Equal(`export VAR1OLD="value1old"
export VAR2OLD="value2old"
export VAR3OLD="value3old"`))
			})
		})

		Context("failure cases", func() {
			It("returns an error when jobs array in pipeline is empty", func() {
				_, err := pc.EnvVars("jobs: []", "some-job", "some-task")
				Expect(err).To(MatchError("pipeline does not contain any jobs - are you sure the pipeline exists?"))
			})

			It("returns error when the given job is not found in the pipeline", func() {
				_, err := pc.EnvVars(modernPipelineYAML, "unknown-job", "some-task")
				Expect(err).To(MatchError(`could not find job "unknown-job" in pipeline`))
			})

			It("returns error when the given task is not found in the pipeline", func() {
				_, err := pc.EnvVars(modernPipelineYAML, "some-job", "unknown-task")
				Expect(err).To(MatchError(`could not find task "unknown-task" in job "some-job" in pipeline`))
			})

			It("returns error when pipeline cannot be unmarshalled", func() {
				_, err := pc.EnvVars("%%%not-valid-yaml%%%", "", "")
				Expect(err).To(MatchError("yaml: could not find expected directive name"))
			})
		})
	})
})
