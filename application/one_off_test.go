package application_test

import (
	"bytes"
	"errors"

	"github.com/kkallday/one-off/application"
	"github.com/kkallday/one-off/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ErrBuffer struct{}

func (*ErrBuffer) Write(_ []byte) (int, error) {
	return -1, errors.New("failed to write")
}

var _ = Describe("one off", func() {
	var (
		fakeFly               *fakes.Fly
		fakePipelineConverter *fakes.PipelineConverter
		stdout                *bytes.Buffer

		oneOff application.OneOff
	)

	BeforeEach(func() {
		fakeFly = &fakes.Fly{}
		fakePipelineConverter = &fakes.PipelineConverter{}
		stdout = &bytes.Buffer{}
		oneOff = application.NewOneOff(fakeFly, fakePipelineConverter, stdout)
	})

	AfterEach(func() {
		application.ResetLookPath()
	})

	It("gets path to fly cli from $PATH", func() {
		application.SetLookPath(func(_ string) (string, error) {
			return "/home/user/bin/fly-cli", nil
		})

		err := oneOff.Run(application.OneOffInputs{})
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeFly.SetPathToFlyCall.CallCount).To(Equal(1))
		Expect(fakeFly.SetPathToFlyCall.Receives.PathToFly).To(Equal("/home/user/bin/fly-cli"))
	})

	Context("when fly override is supplied", func() {
		It("sets the path to fly", func() {
			err := oneOff.Run(application.OneOffInputs{
				FlyOverride: "/some/path/to/fly",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeFly.SetPathToFlyCall.CallCount).To(Equal(1))
			Expect(fakeFly.SetPathToFlyCall.Receives.PathToFly).To(Equal("/some/path/to/fly"))
		})

		It("writes script to stdout with name of program fly override points to", func() {
			fakePipelineConverter.EnvVarsCall.Returns.EnvVars = `export VAR1="foo"`

			err := oneOff.Run(application.OneOffInputs{
				TargetAlias: "some-target-alias",
				Pipeline:    "some-pipeline",
				Job:         "some-job",
				Task:        "some-task",
				FlyOverride: "/some/path/to/custom-fly-cli-program",
			})
			Expect(err).NotTo(HaveOccurred())

			expectedScript := `#!/bin/bash -exu
export VAR1="foo"

custom-fly-cli-program -t some-target-alias execute --config=REPLACE/ME/PATH/TO/TASK --inputs-from some-pipeline/some-job`

			Expect(string(stdout.Bytes())).To(Equal(string(expectedScript)))
		})
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

	It("writes script to stdout", func() {
		fakePipelineConverter.EnvVarsCall.Returns.EnvVars = `export VAR1="foo"
export VAR2="bar"
export VAR3="something else"`
		err := oneOff.Run(application.OneOffInputs{
			TargetAlias: "some-target-alias",
			Pipeline:    "some-pipeline",
			Job:         "some-job",
			Task:        "some-task",
		})
		Expect(err).NotTo(HaveOccurred())

		expectedScript := `#!/bin/bash -exu
export VAR1="foo"
export VAR2="bar"
export VAR3="something else"

fly -t some-target-alias execute --config=REPLACE/ME/PATH/TO/TASK --inputs-from some-pipeline/some-job`

		Expect(string(stdout.Bytes())).To(Equal(string(expectedScript)))
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

		It("returns an error when script cannot be written to stdout", func() {
			errStdout := &ErrBuffer{}
			oneOff = application.NewOneOff(fakeFly, fakePipelineConverter, errStdout)
			err := oneOff.Run(application.OneOffInputs{})
			Expect(err).To(MatchError("failed to write one-off to stdout: failed to write"))
		})
	})
})
