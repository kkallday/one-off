package application_test

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/kkallday/one-off/application"
	"github.com/kkallday/one-off/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("one off", func() {
	var (
		fakeFly               *fakes.Fly
		fakePipelineConverter *fakes.PipelineConverter

		oneOff application.OneOff
	)

	BeforeEach(func() {
		fakeFly = &fakes.Fly{}
		fakePipelineConverter = &fakes.PipelineConverter{}
		oneOff = application.NewOneOff(fakeFly, fakePipelineConverter)
	})

	It("gets pipeline using fly", func() {
		err := oneOff.Run(application.OneOffInputs{
			TargetAlias: "some-target-alias",
			Pipeline:    "some-pipeline",
			Job:         "some-job",
			Task:        "some-task",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeFly.GetPipelineCall.CallCount).To(Equal(1))
		Expect(fakeFly.GetPipelineCall.Receives.TargetAlias).To(Equal("some-target-alias"))
		Expect(fakeFly.GetPipelineCall.Receives.Pipeline).To(Equal("some-pipeline"))
	})

	It("converts the pipeline to env vars", func() {
		fakeFly.GetPipelineCall.Returns.PipelineYAML = "some-pipeline-yaml"
		err := oneOff.Run(application.OneOffInputs{
			TargetAlias: "some-target-alias",
			Pipeline:    "some-pipeline",
			Job:         "some-job",
			Task:        "some-task",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakePipelineConverter.EnvVarsCall.CallCount).To(Equal(1))
		Expect(fakePipelineConverter.EnvVarsCall.Receives.PipelineYAML).To(Equal("some-pipeline-yaml"))
		Expect(fakePipelineConverter.EnvVarsCall.Receives.Job).To(Equal("some-job"))
		Expect(fakePipelineConverter.EnvVarsCall.Receives.Task).To(Equal("some-task"))
	})

	It("writes env vars and fly script to a file", func() {
		fakePipelineConverter.EnvVarsCall.Returns.EnvVars = `export VAR1="foo"
export VAR2="bar"
export VAR3="something else"`
		tempDir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		err = oneOff.Run(application.OneOffInputs{
			TargetAlias: "some-target-alias",
			Pipeline:    "some-pipeline",
			Job:         "some-job",
			Task:        "some-task",
			OutputDir:   tempDir,
		})
		Expect(err).NotTo(HaveOccurred())

		actualScript, err := ioutil.ReadFile(filepath.Join(tempDir, "some-task-one-off"))
		Expect(err).NotTo(HaveOccurred())

		expectedScript := `#!/bin/bash -exu
export VAR1="foo"
export VAR2="bar"
export VAR3="something else"

fly -t some-target-alias execute --config=REPLACE/ME/PATH/TO/TASK --inputs-from some-pipeline/some-job`

		Expect(string(actualScript)).To(Equal(string(expectedScript)))
	})

	Context("failure cases", func() {
		It("returns an error when fly fails to get pipeline", func() {
			fakeFly.GetPipelineCall.Returns.Error = errors.New("some error")
			err := oneOff.Run(application.OneOffInputs{})
			Expect(err).To(MatchError("failed to get pipeline: some error"))
		})

		It("returns an error when pipeline converter fails", func() {
			fakePipelineConverter.EnvVarsCall.Returns.Error = errors.New("failed to convert pipeline")
			err := oneOff.Run(application.OneOffInputs{})
			Expect(err).To(MatchError("failed to retrieve pipeline params from pipeline: failed to convert pipeline"))
		})

		It("returns an error when script file cannot be written", func() {
			err := oneOff.Run(application.OneOffInputs{
				OutputDir: "/some/non/existent/dir",
			})
			Expect(err).To(MatchError("failed to write script: open /some/non/existent/dir/-one-off: no such file or directory"))
		})
	})
})
