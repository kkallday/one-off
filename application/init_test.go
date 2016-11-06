package application_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOneOff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "application")
}
