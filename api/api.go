package api

import (
	"log"
	"net/http"
	"os/exec"
	"fmt"
)

func BastilleCommand(args ...string) (string, error) {

	cmd := exec.Command("bastille", args...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		return output, fmt.Errorf("Bastille %v failed: %v\n%s", args, err, output)
	}

	return output, nil

}

func BastilleCommandLive(args ...string) (string, error) {

	ttydArgs := []string{
		"-t", "disableLeaveAlert=true",
		"-o",
		"-m", "1",
		"-p", "7681",
		"-W",
	}

	var cmdArgs []string
	cmdArgs = append(cmdArgs, ttydArgs...)
	cmdArgs = append(cmdArgs, "bastille")
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("ttyd", cmdArgs...)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf(
			"Bastille live %v failed to start ttyd: %w",
			args,
			err,
		)
	}

	port := fmt.Sprintf("%d", 7681)
	return port, nil
}

func Start() {

	var bindAddr string
	config := loadConfig()
	setConfig(config)

	if Host == "0.0.0.0" || Host == "localhost" || Host == "" {
		bindAddr = "0.0.0.0"
		Host = "localhost"
	} else {
	       bindAddr = Host
	}
	
	addr := fmt.Sprintf("%s:%s", bindAddr, Port)


	loadRoutes()

	log.Println("Starting BastilleBSD API server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}