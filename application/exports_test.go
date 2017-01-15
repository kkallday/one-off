package application

import "os/exec"

func SetLookPath(f func(string) (string, error)) {
	lookPath = f
}

func ResetLookPath() {
	lookPath = exec.LookPath
}
