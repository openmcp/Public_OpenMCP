package util

import "os/exec"

func CmdExec(cmdStr string) (string, error) {

	cmd := exec.Command("bash", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}