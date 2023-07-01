package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func ExecuteBinary(binaryPath string, args []string, v interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf(
			"failed to execute binary %s %v with error: %s, stdout: %s: %w",
			binaryPath,
			args,
			stErr.String(),
			stdOut.String(),
			err,
		)
	}

	if v == nil {
		return bytes.Trim(stdOut.Bytes(), "\n"), nil
	}

	if err := json.Unmarshal(stdOut.Bytes(), v); err != nil {
		// TODO Maybe failback to text parsing instead??
		return nil, fmt.Errorf(
			"failed to unmarshal binary response to given struct(%t): %w",
			v,
			err,
		)
	}

	return nil, nil
}
