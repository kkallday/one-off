package application_test

import (
	"errors"

	"github.com/kkallday/one-off/application"
	"github.com/kkallday/one-off/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("app", func() {
	var (
		fakeOneOff *fakes.OneOff

		app application.App
	)

	BeforeEach(func() {
		fakeOneOff = &fakes.OneOff{}
		app = application.New(fakeOneOff)
	})

	It("runs one-off", func() {
		err := app.Execute([]string{
			"-ta", "some-target-alias",
			"-p", "some-pipeline",
			"-j", "some-job",
			"-t", "some-task",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeOneOff.RunCall.CallCount).To(Equal(1))
		Expect(fakeOneOff.RunCall.Receives.OneOffInputs).To(Equal(application.OneOffInputs{
			TargetAlias: "some-target-alias",
			Pipeline:    "some-pipeline",
			Job:         "some-job",
			Task:        "some-task",
		}))
	})

	Context("failure cases", func() {
		It("returns an error when undefined args are provided", func() {
			err := app.Execute([]string{
				"-ta", "some-target-alias",
				"-p", "some-pipeline",
				"-j", "some-job",
				"-t", "some-task",
				"-foo", "bar",
			})
			Expect(err).To(MatchError("flag provided but not defined: -foo"))
		})

		DescribeTable("required args", func(args []string, expectedErr string) {
			err := app.Execute(args)
			Expect(err).To(MatchError(expectedErr))
		},
			Entry("missing target alias", []string{
				"-p", "some-pipeline",
				"-j", "some-job",
				"-t", "some-task",
			}, "missing target alias -ta"),
			Entry("missing pipeline", []string{
				"-ta", "some-target-alias",
				"-j", "some-job",
				"-t", "some-task",
			}, "missing pipeline -p"),
			Entry("missing job", []string{
				"-ta", "some-target-alias",
				"-p", "some-pipeline",
				"-t", "some-task",
			}, "missing job -j"),
			Entry("missing task", []string{
				"-ta", "some-target-alias",
				"-p", "some-pipeline",
				"-j", "some-job",
			}, "missing task -t"),
			Entry("missing all args", []string{}, "missing target alias -ta, pipeline -p, job -j, task -t"),
		)

		It("returns an error when one off fails", func() {
			fakeOneOff.RunCall.Returns.Error = errors.New("run failed")
			err := app.Execute([]string{
				"-ta", "some-target-alias",
				"-p", "some-pipeline",
				"-j", "some-job",
				"-t", "some-task",
			})
			Expect(err).To(MatchError("run failed"))
		})

	})
})
