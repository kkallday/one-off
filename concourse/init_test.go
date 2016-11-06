package concourse_test

import (
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConcourse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "concourse")
}

var (
	pathToFakeFly        string
	pathToFakeErroredFly string
)

var _ = BeforeSuite(func() {
	var err error
	pathToFakeFly, err = gexec.Build("github.com/kkallday/one-off/fakes/fly")
	Expect(err).NotTo(HaveOccurred())

	pathToFakeErroredFly, err = gexec.Build("github.com/kkallday/one-off/fakes/erroredfly")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
