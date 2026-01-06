package api

import (
	"fmt"
	"os/exec"
)

// Runs the main Bastille command after everything has been validated
func BastilleCommand(args ...string) (string, error) {

	cmd := exec.Command("bastille", args...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		return output, fmt.Errorf("Bastille %v failed: %v\n%s", args, err, output)
	}

	return output, nil

}
