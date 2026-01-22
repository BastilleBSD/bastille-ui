package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func RocinanteCommand(args ...string) (string, error) {

	cmd := exec.Command("rocinante", args...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		return output, fmt.Errorf("Rocinante %v failed: %v\n%s", args, err, output)
	}

	return output, nil
}

func RocinanteCommandLive(args ...string) (string, error) {

	ttydArgs := []string{
		"-t", "disableLeaveAlert=true",
		"-o",
		"--ipv6",
		"-m", "1",
		"-p", "7681",
		"-W",
	}

	var cmdArgs []string
	cmdArgs = append(cmdArgs, ttydArgs...)
	cmdArgs = append(cmdArgs, "rocinante")
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("ttyd", cmdArgs...)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf(
			"Rocinante live %v failed to start ttyd: %w",
			args,
			err,
		)
	}

	port := fmt.Sprintf("%d", 7681)
	return port, nil
}

func ParseAndRunRocinanteCommand(w http.ResponseWriter, r *http.Request, cmdArgs []string) {

	if err := ValidateRocinanteCommandParameters(r, cmdArgs); err != nil {
		logAll("error", r, cmdArgs, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isLive := strings.Contains(r.URL.Path, "/api/v1/rocinante/live/")

	var (
		result RocinanteCommandOutputStruct
		err    error
	)

	if isLive {
		result.port, err = RocinanteCommandLive(cmdArgs...)
	} else {
		result.output, err = RocinanteCommand(cmdArgs...)
	}

	if err != nil {
		logAll("error", r, cmdArgs, fmt.Sprintf("failed: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isLive {
		w.Header().Set("X-TTYD-Port", result.port)
	} else {
		fmt.Fprint(w, result.output)
	}

	logAll("info", r, cmdArgs, "success")
}

func RocinanteBootstrapHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"bootstrap"}

	options := getParam(r, "options")
	url := getParam(r, "url")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if url == "" {
		http.Error(w, "[ERROR]: Missing url parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, url)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteCmdHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"cmd"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteLimitsHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"limits"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteListHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"list"}

	options := getParam(r, "options")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinantePkgHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"pkg"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteServiceHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"service"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteSysctlHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"sysctl"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteSysrcHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"sysrc"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteTemplateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"template"}

	options := getParam(r, "options")
	action := getParam(r, "action")
	template := getParam(r, "template")
	args := getParam(r, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if action != "" {
		cmdArgs = append(cmdArgs, action)
	}
	if template == "" {
		http.Error(w, "[ERROR]: Missing template parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, template)
	if args != "" {
		cmdArgs = append(cmdArgs, strings.Fields(args)...)
	}

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteUpdateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"update"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteUpgradeHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"upgrade"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteVerifyHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"verify"}

	options := getParam(r, "options")
	template := getParam(r, "template")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if template == "" {
		http.Error(w, "[ERROR]: Missing template parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, template)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteZfsHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"zfs"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}

func RocinanteZpoolHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"zpool"}

	args := getParam(r, "args")

	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(w, r, cmdArgs)
}
