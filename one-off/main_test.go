package main_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("one off", func() {
	It("prints pipeline definition to stdout", func() {
		args := []string{
			"-ta", "some-target-alias",
			"-p", "some-pipeline",
			"-j", "some-job",
			"-t", "some-task",
			"-fo", pathToFakeFly,
		}

		cmd := exec.Command(pathToOneOff, args...)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		buf, err := ioutil.ReadFile(filepath.Join("fixtures", "basic_script"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(session.Out.Contents())).To(Equal(strings.Trim(string(buf), "\n")))
	})
})
