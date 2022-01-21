package shell

import (
	"os"
	"os/exec"
)

func convertOutput(output []byte, err error) (string, error) {
	return string(output), err
}

func Exec(cmd string, params ...string) (string, error) {
	return convertOutput(exec.Command(cmd, params...).CombinedOutput())
}

func ExecWithEnv(env []string, cmd string, params ...string) (string, error) {
	command := exec.Command(cmd, params...)
	command.Env = os.Environ()
	for _, envVar := range env {
		command.Env = append(command.Env, envVar)
	}

	return convertOutput(command.CombinedOutput())
}
