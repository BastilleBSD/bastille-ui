package api

import (
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var bastilleLock sync.Mutex

func BastilleCommand(args ...string) (string, error) {

	logRequest("debug", "BastilleCommand", nil, args, nil)

	bastilleLock.Lock()
	defer bastilleLock.Unlock()

	cmd := exec.Command("bastille", args...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		logRequest("error", "command failed", nil, args, output)
		return "", err
	}

	return output, nil
}

func BastilleCommandLive(args ...string) (string, error) {

	logRequest("debug", "BastilleCommandLive", nil, args, nil)

	bastilleLock.Lock()
	defer bastilleLock.Unlock()

	ttydArgs := []string{
		"-t", "disableLeaveAlert=true",
		"-o",
		"--ipv6",
		"-m", "1",
		"-p", "7681",
		"-W",
	}

	cmdArgs := append(ttydArgs, "bastille")
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("ttyd", cmdArgs...)
	if err := cmd.Start(); err != nil {
		logRequest("error", "ttyd command failed", nil, args, err)
		return "", err
	}

	return "7681", nil
}

func ParseAndRunBastilleCommand(c *gin.Context, cmdArgs []string) {

	logRequest("debug", "ParseAndRunBastilleCommand", c, cmdArgs, nil)

	if err := ValidateBastilleCommandParameters(c, cmdArgs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logRequest("error", "parameter validation failed", c, cmdArgs, err)
		return
	}

	isLive := strings.Contains(c.FullPath(), "/api/v1/bastille/live/")
	var result string
	var err error

	if isLive {
		result, err = BastilleCommandLive(cmdArgs...)
	} else {
		result, err = BastilleCommand(cmdArgs...)
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

func BastilleBootstrapHandler(c *gin.Context) {

	logRequest("debug", "BastilleBootstrapHandler", c, nil, nil)

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

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleCloneHandler(c *gin.Context) {

	logRequest("debug", "BastilleCloneHandler", c, nil, nil)

	cmdArgs := []string{"clone"}
	options := getParam(c, "options")
	target := getParam(c, "target")
	new_name := getParam(c, "new_name")
	ip := getParam(c, "ip")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing target parameter"})
		logRequest("error", "missing target parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if new_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing new_name parameter"})
		logRequest("error", "missing new_name parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, new_name)

	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ip parameter"})
		logRequest("error", "missing ip parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, ip)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleCmdHandler(c *gin.Context) {

	logRequest("debug", "BastilleCmdHandler", c, nil, nil)

	cmdArgs := []string{"cmd"}
	options := getParam(c, "options")
	target := getParam(c, "target")
	command := getParam(c, "command")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing target parameter"})
		logRequest("error", "missing target parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing command parameter"})
		logRequest("error", "missing command parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(command)...)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleConfigHandler(c *gin.Context) {

	logRequest("debug", "BastilleConfigHandler", c, nil, nil)

	cmdArgs := []string{"config"}
	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	property := getParam(c, "property")
	value := getParam(c, "value")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing target parameter"})
		logRequest("error", "missing target parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if action != "set" && action != "add" && action != "get" && action != "remove" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown action parameter"})
		logRequest("error", "unknown action parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, action)

	if property == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing property parameter"})
		logRequest("error", "missing property parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, property)

	if value != "" {
		cmdArgs = append(cmdArgs, value)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleConsoleHandler(c *gin.Context) {

	logRequest("debug", "BastilleConsoleHandler", c, nil, nil)

	cmdArgs := []string{"console"}
	options := getParam(c, "options")
	target := getParam(c, "target")
	user := getParam(c, "user")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing target parameter"})
		logRequest("error", "missing target parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if user != "" {
		cmdArgs = append(cmdArgs, user)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleConvertHandler(c *gin.Context) {

	logRequest("debug", "BastilleConvertHandler", c, nil, nil)

	cmdArgs := []string{"convert"}
	options := getParam(c, "options")
	target := getParam(c, "target")
	release := getParam(c, "release")

	if options != "" {
		options = options + " -ay"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	} else {
		cmdArgs = append(cmdArgs, "-ay")
	}

	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing target parameter"})
		logRequest("error", "missing target parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if release != "" {
		cmdArgs = append(cmdArgs, release)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleCpHandler(c *gin.Context) {

	logRequest("debug", "BastilleCpHandler", c, nil, nil)

	cmdArgs := []string{"cp"}
	options := getParam(c, "options")
	target := getParam(c, "target")
	host_path := getParam(c, "host_path")
	jail_path := getParam(c, "jail_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing target parameter"})
		logRequest("error", "missing target parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if host_path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing host_path parameter"})
		logRequest("error", "missing host_path parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, host_path)

	if jail_path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing jail_path parameter"})
		logRequest("error", "missing jail_path parameter", c, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, jail_path)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleCreateHandler(c *gin.Context) {

	logRequest("debug", "BastilleCreateHandler", nil, nil, nil)

	cmdArgs := []string{"create"}

	options := getParam(c, "options")
	name := getParam(c, "name")
	release := getParam(c, "release")
	ip := getParam(c, "ip")
	iface := getParam(c, "iface")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, "Missing name parameter.")
		logRequest("error", "missing name parameter", nil, cmdArgs, nil)
		return
	}
	if release == "" {
		c.JSON(http.StatusBadRequest, "Missing release parameter")
		logRequest("error", "missing release parameter", nil, cmdArgs, nil)
		return
	}
	if ip == "" {
		c.JSON(http.StatusBadRequest, "Missing ip parameter")
		logRequest("error", "missing ip parameter", nil, cmdArgs, nil)
		return
	}

	cmdArgs = append(cmdArgs, name, release, ip)

	if iface != "" {
		cmdArgs = append(cmdArgs, iface)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleDestroyHandler(c *gin.Context) {

	logRequest("debug", "BastilleDestroyHandler", nil, nil, nil)

	cmdArgs := []string{"destroy"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleEditHandler(c *gin.Context) {

	logRequest("debug", "BastilleEditHandler", nil, nil, nil)

	cmdArgs := []string{"edit"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	file := getParam(c, "file")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if file != "" {
		cmdArgs = append(cmdArgs, file)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleEtcupdateHandler(c *gin.Context) {

	logRequest("debug", "BastilleEtcupdateHandler", nil, nil, nil)

	cmdArgs := []string{"etcupdate"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	release := getParam(c, "release")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}

	if action == "bootstrap" {
		if release == "" {
			c.JSON(http.StatusBadRequest, "Missing release parameter")
			logRequest("error", "missing release parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, action, release)
	} else {
		if target == "" {
			c.JSON(http.StatusBadRequest, "Missing target parameter")
			logRequest("error", "missing target parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, target)
		if action == "update" {
			if release == "" {
				c.JSON(http.StatusBadRequest, "Missing release parameter")
				logRequest("error", "missing release parameter", nil, cmdArgs, nil)
				return
			}
			cmdArgs = append(cmdArgs, release)
		} else {
			if action == "" {
				c.JSON(http.StatusBadRequest, "Missing action parameter")
				logRequest("error", "missing action parameter", nil, cmdArgs, nil)
				return
			}
			cmdArgs = append(cmdArgs, action)
		}
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleExportHandler(c *gin.Context) {

	logRequest("debug", "BastilleExportHandler", nil, nil, nil)

	cmdArgs := []string{"export"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	path := getParam(c, "path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if path != "" {
		cmdArgs = append(cmdArgs, path)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleHtopHandler(c *gin.Context) {

	logRequest("debug", "BastilleHtopHandler", nil, nil, nil)

	cmdArgs := []string{"htop"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleImportHandler(c *gin.Context) {

	logRequest("debug", "BastilleImportHandler", nil, nil, nil)

	cmdArgs := []string{"import"}

	options := getParam(c, "options")
	file := getParam(c, "file")
	release := getParam(c, "release")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if file == "" {
		c.JSON(http.StatusBadRequest, "Missing file parameter")
		logRequest("error", "missing file parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, file)
	if release != "" {
		cmdArgs = append(cmdArgs, release)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleJcpHandler(c *gin.Context) {

	logRequest("debug", "BastilleJcpHandler", nil, nil, nil)

	cmdArgs := []string{"jcp"}

	options := getParam(c, "options")
	source_jail := getParam(c, "source_jail")
	source_path := getParam(c, "source_path")
	destination_jail := getParam(c, "destination_jail")
	destination_path := getParam(c, "destination_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if source_jail == "" {
		c.JSON(http.StatusBadRequest, "Missing source_jail parameter")
		logRequest("error", "missing source_jail parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, source_jail)
	if source_path == "" {
		c.JSON(http.StatusBadRequest, "Missing source_path parameter")
		logRequest("error", "missing source_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, source_path)
	if destination_jail == "" {
		c.JSON(http.StatusBadRequest, "Missing destination_jail parameter")
		logRequest("error", "missing destination_jail parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, destination_jail)
	if destination_path == "" {
		c.JSON(http.StatusBadRequest, "Missing destination_path parameter")
		logRequest("error", "missing destination_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, destination_path)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleLimitsHandler(c *gin.Context) {

	logRequest("debug", "BastilleLimitsHandler", nil, nil, nil)

	cmdArgs := []string{"limits"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	args := getParam(c, "args")
	option := getParam(c, "option")
	value := getParam(c, "value")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "" {
		c.JSON(http.StatusBadRequest, "Missing action parameter")
		logRequest("error", "missing action parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, action)

	switch action {
	case "add":
		if option == "" {
			c.JSON(http.StatusBadRequest, "Missing option parameter")
			logRequest("error", "missing option parameter", nil, cmdArgs, nil)
			return
		}
		if value == "" {
			c.JSON(http.StatusBadRequest, "Missing value parameter")
			logRequest("error", "missing value parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, option, value)
	case "remove":
		if option == "" {
			c.JSON(http.StatusBadRequest, "Missing option parameter")
			logRequest("error", "missing option parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, option)
	case "clear", "reset", "stats":
		// just append the action
	case "list", "show":
		if args == "active" {
			cmdArgs = append(cmdArgs, action, args)
		} else {
			cmdArgs = append(cmdArgs, action)
		}
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleListHandler(c *gin.Context) {

	logRequest("debug", "BastilleListHandler", nil, nil, nil)

	cmdArgs := []string{"list"}

	options := getParam(c, "options")
	item := getParam(c, "item")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if item != "" {
		cmdArgs = append(cmdArgs, item)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleMigrateHandler(c *gin.Context) {

	logRequest("debug", "BastilleMigrateHandler", nil, nil, nil)

	cmdArgs := []string{"migrate"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	destination := getParam(c, "destination")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if destination == "" {
		c.JSON(http.StatusBadRequest, "Missing destination parameter")
		logRequest("error", "missing destination parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, destination)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleMonitorHandler(c *gin.Context) {

	logRequest("debug", "BastilleMonitorHandler", nil, nil, nil)

	cmdArgs := []string{"monitor"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	service := getParam(c, "service")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}

	if action == "enable" || action == "disable" || action == "status" {
		cmdArgs = append(cmdArgs, action)
	} else if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	} else {
		cmdArgs = append(cmdArgs, target)
	}

	if action == "add" || action == "delete" {
		if service == "" {
			c.JSON(http.StatusBadRequest, "Missing service parameter")
			logRequest("error", "missing service parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, action, service)
	} else if action == "list" {
		cmdArgs = append(cmdArgs, action)
		if service != "" {
			cmdArgs = append(cmdArgs, service)
		}
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleMountHandler(c *gin.Context) {

	logRequest("debug", "BastilleMountHandler", nil, nil, nil)

	cmdArgs := []string{"mount"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	host_path := getParam(c, "host_path")
	jail_path := getParam(c, "jail_path")
	fs_type := getParam(c, "fs_type")
	fs_options := getParam(c, "fs_options")
	dump := getParam(c, "dump")
	pass_number := getParam(c, "pass_number")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if host_path == "" {
		c.JSON(http.StatusBadRequest, "Missing host_path parameter")
		logRequest("error", "missing host_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, host_path)
	if jail_path == "" {
		c.JSON(http.StatusBadRequest, "Missing jail_path parameter")
		logRequest("error", "missing jail_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, jail_path)

	if fs_type != "" || fs_options != "" || dump != "" || pass_number != "" {
		if fs_type == "" || fs_options == "" || dump == "" || pass_number == "" {
			c.JSON(http.StatusBadRequest, "Missing mount parameter(s)")
			logRequest("error", "missing mount parameter(s)", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, fs_type, fs_options, dump, pass_number)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleNetworkHandler(c *gin.Context) {

	logRequest("debug", "BastilleNetworkHandler", nil, nil, nil)

	cmdArgs := []string{"network"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	iface := getParam(c, "iface")
	ip := getParam(c, "ip")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if action == "add" {
		if iface == "" {
			c.JSON(http.StatusBadRequest, "Missing iface parameter")
			logRequest("error", "missing iface parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, action, iface)
		if ip != "" {
			cmdArgs = append(cmdArgs, ip)
		}
	} else if action == "remove" {
		if iface == "" {
			c.JSON(http.StatusBadRequest, "Missing iface parameter")
			logRequest("error", "missing iface parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, action, iface)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastillePkgHandler(c *gin.Context) {

	logRequest("debug", "BastillePkgHandler", nil, nil, nil)

	cmdArgs := []string{"pkg"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	args := getParam(c, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if args == "" {
		c.JSON(http.StatusBadRequest, "Missing args parameter")
		logRequest("error", "missing args parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleRcpHandler(c *gin.Context) {

	logRequest("debug", "BastilleRcpHandler", nil, nil, nil)

	cmdArgs := []string{"rcp"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	jailPath := getParam(c, "jail_path")
	hostPath := getParam(c, "host_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if jailPath == "" {
		c.JSON(http.StatusBadRequest, "Missing jail_path parameter")
		logRequest("error", "missing jail_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, jailPath)

	if hostPath == "" {
		c.JSON(http.StatusBadRequest, "Missing host_path parameter")
		logRequest("error", "missing host_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, hostPath)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleRdrHandler(c *gin.Context) {

	logRequest("debug", "BastilleRdrHandler", nil, nil, nil)

	cmdArgs := []string{"rdr"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	protocol := getParam(c, "protocol")
	hostPort := getParam(c, "host_port")
	jailPort := getParam(c, "jail_port")
	logOptions := getParam(c, "log_options")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if action == "clear" || action == "reset" || action == "list" {
		cmdArgs = append(cmdArgs, action)
	} else {
		if protocol == "" {
			c.JSON(http.StatusBadRequest, "Missing protocol parameter")
			logRequest("error", "missing protocol parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, protocol)

		if hostPort == "" {
			c.JSON(http.StatusBadRequest, "Missing host_port parameter")
			logRequest("error", "missing host_port parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, hostPort)

		if jailPort == "" {
			c.JSON(http.StatusBadRequest, "Missing jail_port parameter")
			logRequest("error", "missing jail_port parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, jailPort)

		if action == "log" {
			if logOptions == "" {
				c.JSON(http.StatusBadRequest, "Missing log_options parameter")
				logRequest("error", "missing log_options parameter", nil, cmdArgs, nil)
				return
			}
			cmdArgs = append(cmdArgs, action, logOptions)
		}
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleRenameHandler(c *gin.Context) {

	logRequest("debug", "BastilleRenameHandler", nil, nil, nil)

	cmdArgs := []string{"rename"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	newName := getParam(c, "new_name")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	if newName == "" {
		c.JSON(http.StatusBadRequest, "Missing new_name parameter")
		logRequest("error", "missing new_name parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target, newName)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleRestartHandler(c *gin.Context) {

	logRequest("debug", "BastilleRestartHandler", nil, nil, nil)

	cmdArgs := []string{"restart"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleServiceHandler(c *gin.Context) {

	logRequest("debug", "BastilleServiceHandler", nil, nil, nil)

	cmdArgs := []string{"service"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	service := getParam(c, "service")
	args := getParam(c, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if service == "" {
		c.JSON(http.StatusBadRequest, "Missing service parameter")
		logRequest("error", "missing service parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, service)
	if args == "" {
		c.JSON(http.StatusBadRequest, "Missing args parameter")
		logRequest("error", "missing args parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleSetupHandler(c *gin.Context) {

	logRequest("debug", "BastilleSetupHandler", nil, nil, nil)

	cmdArgs := []string{"setup"}

	options := getParam(c, "options")
	item := getParam(c, "item")
	args := getParam(c, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if item != "" {
		cmdArgs = append(cmdArgs, item)
	}
	if args != "" {
		cmdArgs = append(cmdArgs, args)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleStartHandler(c *gin.Context) {

	logRequest("debug", "BastilleStartHandler", nil, nil, nil)

	cmdArgs := []string{"start"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleStopHandler(c *gin.Context) {

	logRequest("debug", "BastilleStopHandler", nil, nil, nil)

	cmdArgs := []string{"stop"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleSysrcHandler(c *gin.Context) {

	logRequest("debug", "BastilleSysrcHandler", nil, nil, nil)

	cmdArgs := []string{"sysrc"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	args := getParam(c, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if args == "" {
		c.JSON(http.StatusBadRequest, "Missing args parameter")
		logRequest("error", "missing args parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleTagsHandler(c *gin.Context) {

	logRequest("debug", "BastilleTagsHandler", nil, nil, nil)

	cmdArgs := []string{"tags"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	tags := getParam(c, "tags")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if action == "add" || action == "delete" {
		if tags == "" {
			c.JSON(http.StatusBadRequest, "Missing tags parameter")
			logRequest("error", "missing tags parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, action, tags)
	} else if action == "list" {
		cmdArgs = append(cmdArgs, action)
		if tags != "" {
			cmdArgs = append(cmdArgs, tags)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid action parameter")
		logRequest("error", "invalid action parameter", nil, cmdArgs, nil)
		return
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleTemplateHandler(c *gin.Context) {

	logRequest("debug", "BastilleTemplateHandler", nil, nil, nil)

	cmdArgs := []string{"template"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	template := getParam(c, "template")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}

	if action == "convert" {
		cmdArgs = append(cmdArgs, action)
		if template == "" {
			c.JSON(http.StatusBadRequest, "Missing template parameter")
			logRequest("error", "missing template parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, template)
	} else {
		if target == "" {
			c.JSON(http.StatusBadRequest, "Missing target parameter")
			logRequest("error", "missing target parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, target)
		if template == "" {
			c.JSON(http.StatusBadRequest, "Missing template parameter")
			logRequest("error", "missing template parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, template)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleTopHandler(c *gin.Context) {

	logRequest("debug", "BastilleTopHandler", nil, nil, nil)

	cmdArgs := []string{"top"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleUmountHandler(c *gin.Context) {

	logRequest("debug", "BastilleUmountHandler", nil, nil, nil)

	cmdArgs := []string{"umount"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	jailPath := getParam(c, "jail_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if jailPath == "" {
		c.JSON(http.StatusBadRequest, "Missing jail_path parameter")
		logRequest("error", "missing jail_path parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, jailPath)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleUpdateHandler(c *gin.Context) {

	logRequest("debug", "BastilleUpdateHandler", nil, nil, nil)

	cmdArgs := []string{"update"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleUpgradeHandler(c *gin.Context) {

	logRequest("debug", "BastilleUpgradeHandler", nil, nil, nil)

	cmdArgs := []string{"upgrade"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	newRelease := getParam(c, "new_release")
	action := getParam(c, "action")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	if action == "install" {
		cmdArgs = append(cmdArgs, action)
	} else {
		if newRelease == "" {
			c.JSON(http.StatusBadRequest, "Missing new_release parameter")
			logRequest("error", "missing new_release parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, newRelease)
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleVerifyHandler(c *gin.Context) {

	logRequest("debug", "BastilleVerifyHandler", nil, nil, nil)

	cmdArgs := []string{"verify"}

	options := getParam(c, "options")
	target := getParam(c, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunBastilleCommand(c, cmdArgs)
}

func BastilleZfsHandler(c *gin.Context) {

	logRequest("debug", "BastilleZfsHandler", nil, nil, nil)

	cmdArgs := []string{"zfs"}

	options := getParam(c, "options")
	target := getParam(c, "target")
	action := getParam(c, "action")
	tag := getParam(c, "tag")
	keyValue := getParam(c, "key_value")
	dataset := getParam(c, "dataset")
	jailPath := getParam(c, "jail_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		c.JSON(http.StatusBadRequest, "Missing target parameter")
		logRequest("error", "missing target parameter", nil, cmdArgs, nil)
		return
	}
	cmdArgs = append(cmdArgs, target)

	switch action {
	case "snapshot", "destroy", "rollback":
		cmdArgs = append(cmdArgs, action)
		if tag != "" {
			cmdArgs = append(cmdArgs, tag)
		}
	case "df", "usage":
		cmdArgs = append(cmdArgs, action)
	case "get", "set":
		cmdArgs = append(cmdArgs, action)
		if keyValue == "" {
			c.JSON(http.StatusBadRequest, "Missing key_value parameter")
			logRequest("error", "missing key_value parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, keyValue)
	case "jail":
		cmdArgs = append(cmdArgs, action)
		if dataset == "" {
			c.JSON(http.StatusBadRequest, "Missing dataset parameter")
			logRequest("error", "missing dataset parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, dataset)
		if jailPath == "" {
			c.JSON(http.StatusBadRequest, "Missing jail_path parameter")
			logRequest("error", "missing jail_path parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, jailPath)
	case "unjail":
		cmdArgs = append(cmdArgs, action)
		if jailPath == "" {
			c.JSON(http.StatusBadRequest, "Missing jail_path parameter")
			logRequest("error", "missing jail_path parameter", nil, cmdArgs, nil)
			return
		}
		cmdArgs = append(cmdArgs, jailPath)
	default:
		c.JSON(http.StatusBadRequest, "Invalid action parameter")
		logRequest("error", "invalid action parameter", nil, cmdArgs, nil)
		return
	}

	ParseAndRunBastilleCommand(c, cmdArgs)
}