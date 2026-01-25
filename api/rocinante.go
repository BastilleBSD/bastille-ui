package api

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func RocinanteCommand(args ...string) (string, error) {

	logRequest("debug", "RocinanteCommand", nil, args, nil)

	cmd := exec.Command("rocinante", args...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		logRequest("error", "command failed", nil, args, output)
		return "", err
	}

	return output, nil
}

func RocinanteCommandLive(args ...string) (string, error) {

	logRequest("debug", "RocinanteCommandLive", nil, args, nil)

	ttydArgs := []string{
		"-t", "disableLeaveAlert=true",
		"-o",
		"--ipv6",
		"-m", "1",
		"-p", "7681",
		"-W",
	}

	cmdArgs := append(ttydArgs, "rocinante")
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("ttyd", cmdArgs...)
	if err := cmd.Start(); err != nil {
		logRequest("error", "ttyd command failed", nil, args, err)
		return "", err
	}

	return "7681", nil
}

func ParseAndRunRocinanteCommand(c *gin.Context, cmdArgs []string) {

	logRequest("debug", "ParseAndRunRocinanteCommand", c, cmdArgs, nil)

	if err := ValidateRocinanteCommandParameters(c, cmdArgs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logRequest("error", "parameter validation failed", c, cmdArgs, err)
		return
	}

	isLive := strings.Contains(c.FullPath(), "/api/v1/rocinante/live/")
	var result string
	var err error

	if isLive {
		result, err = RocinanteCommandLive(cmdArgs...)
	} else {
		result, err = RocinanteCommand(cmdArgs...)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logRequest("error", "command failed", c, cmdArgs, err)
		return
	}

	if isLive {
		c.Header("X-TTYD-Port", result)
		c.JSON(http.StatusOK, gin.H{"port": result})
		logRequest("info", "success (live)", c, cmdArgs, result)
	} else {
		c.String(http.StatusOK, result)
		logRequest("info", "success", c, cmdArgs, result)
	}
}

func RocinanteBootstrapHandler(c *gin.Context) {

	logRequest("debug", "RocinanteBootstrapHandler", c, nil, nil)

	cmdArgs := []string{"bootstrap"}
	options := getParam(c, "options")
	url := getParam(c, "url")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing url parameter"})
		logRequest("error", "missing url parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, url)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteCmdHandler(c *gin.Context) {

	logRequest("debug", "RocinanteCmdHandler", c, nil, nil)

	cmdArgs := []string{"cmd"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteLimitsHandler(c *gin.Context) {

	logRequest("debug", "RocinanteLimitsHandler", c, nil, nil)

	cmdArgs := []string{"limits"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteListHandler(c *gin.Context) {

	logRequest("debug", "RocinanteListHandler", c, nil, nil)

	cmdArgs := []string{"list"}
	options := getParam(c, "options")
	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinantePkgHandler(c *gin.Context) {

	logRequest("debug", "RocinantePkgHandler", c, nil, nil)

	cmdArgs := []string{"pkg"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteServiceHandler(c *gin.Context) {

	logRequest("debug", "RocinanteServiceHandler", c, nil, nil)

	cmdArgs := []string{"service"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteSysctlHandler(c *gin.Context) {

	logRequest("debug", "RocinanteSysctlHandler", c, nil, nil)

	cmdArgs := []string{"sysctl"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteSysrcHandler(c *gin.Context) {

	logRequest("debug", "RocinanteSysrcHandler", c, nil, nil)

	cmdArgs := []string{"sysrc"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteTemplateHandler(c *gin.Context) {

	logRequest("debug", "RocinanteTemplateHandler", c, nil, nil)

	cmdArgs := []string{"template"}
	options := getParam(c, "options")
	action := getParam(c, "action")
	template := getParam(c, "template")
	args := getParam(c, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if action != "" {
		cmdArgs = append(cmdArgs, action)
	}
	if template == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template parameter"})
		logRequest("error", "missing template parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, template)
	if args != "" {
		cmdArgs = append(cmdArgs, strings.Fields(args)...)
	}

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteUpdateHandler(c *gin.Context) {

	logRequest("debug", "RocinanteUpdateHandler", c, nil, nil)

	cmdArgs := []string{"update"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteUpgradeHandler(c *gin.Context) {

	logRequest("debug", "RocinanteUpgradeHandler", c, nil, nil)

	cmdArgs := []string{"upgrade"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteVerifyHandler(c *gin.Context) {

	logRequest("debug", "RocinanteVerifyHandler", c, nil, nil)

	cmdArgs := []string{"verify"}
	options := getParam(c, "options")
	template := getParam(c, "template")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if template == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template parameter"})
		logRequest("error", "missing template parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, template)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteZfsHandler(c *gin.Context) {

	logRequest("debug", "RocinanteZfsHandler", c, nil, nil)

	cmdArgs := []string{"zfs"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}

func RocinanteZpoolHandler(c *gin.Context) {

	logRequest("debug", "RocinanteZpoolHandler", c, nil, nil)

	cmdArgs := []string{"zpool"}
	args := getParam(c, "args")
	if args == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing args parameter"})
		logRequest("error", "missing args parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunRocinanteCommand(c, cmdArgs)
}