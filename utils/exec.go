package utils

import (
	"bytes"
	"dlx/log"
	"fmt"
	"os/exec"
)

func ExecuteCmd(name string, args ...string) (int, string, string) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		e := fmt.Errorf("execute command: %s, err: %v, stdout: %s, stderr: %s",
			cmd.String(), err, stdout.String(), stderr.String())
		log.Sugar.Error(e)
	}

	return cmd.ProcessState.ExitCode(), stdout.String(), stderr.String()
}
