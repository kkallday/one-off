package main_test

import (
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOneOff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "one-off")
}

var (
	pathToOneOff  string
	pathToFakeFly string
)

var _ = BeforeSuite(func() {
	var err error
	pathToOneOff, err = gexec.Build("github.com/kkallday/one-off/one-off")
	Expect(err).NotTo(HaveOccurred())

	pathToFakeFly, err = gexec.Build("github.com/kkallday/one-off/one-off/fly")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
