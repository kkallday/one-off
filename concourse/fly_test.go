package concourse_test

import (
	"github.com/kkallday/one-off/concourse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("fly", func() {
	var fly concourse.Fly

	BeforeEach(func() {
		fly = concourse.NewFly()
		fly.SetPathToFly(pathToFakeFly)
	})

	Describe("GetPipeline", func() {
		It("shells out to fly", func() {
			pipeline, err := fly.GetPipeline("some-target", "some-pipeline")
			Expect(err).NotTo(HaveOccurred())

			// fake fly binary echoes args it was called with. tests that the stdout is returned
			Expect(pipeline).To(Equal("-t some-target get-pipeline --pipeline some-pipeline"))
		})

		Context("failure cases", func() {
			It("returns an error when fly fails to run", func() {
				fly = concourse.NewFly()
				fly.SetPathToFly("unknown-command")

				_, err := fly.GetPipeline("", "")
				Expect(err).To(MatchError(`exec: "unknown-command": executable file not found in $PATH`))
			})

			It("returns an error containing stderr when fly exits with a non-zero exit code", func() {
				fly = concourse.NewFly()
				fly.SetPathToFly(pathToFakeErroredFly)

				_, err := fly.GetPipeline("", "")
				Expect(err).To(MatchError("exit status 1\nstderr from fly: some error message"))
			})
		})
	})
})
