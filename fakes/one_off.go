package fakes

import "github.com/kkallday/one-off/application"

type OneOff struct {
	RunCall struct {
		CallCount int
		Receives  struct {
			OneOffInputs application.OneOffInputs
		}
		Returns struct {
			Error error
		}
	}
}

func (o *OneOff) Run(oneOffInputs application.OneOffInputs) error {
	o.RunCall.CallCount++
	o.RunCall.Receives.OneOffInputs = oneOffInputs
	return o.RunCall.Returns.Error
}
