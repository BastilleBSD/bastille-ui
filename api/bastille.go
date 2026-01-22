package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
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
		"--ipv6",
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

func ParseAndRunCommand(w http.ResponseWriter, r *http.Request, cmdArgs []string) {

	if err := ValidateBastilleCommandParameters(r, cmdArgs); err != nil {
		logAll("error", r, cmdArgs, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isLive := strings.Contains(r.URL.Path, "/api/v1/bastille/live/")

	var (
		result BastilleCommandOutputStruct
		err    error
	)

	if isLive {
		result.port, err = BastilleCommandLive(cmdArgs...)
	} else {
		result.output, err = BastilleCommand(cmdArgs...)
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

func BastilleBootstrapHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"bootstrap"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	arch := getParam(r, "arch")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if arch != "" {
		cmdArgs = append(cmdArgs, arch)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleCloneHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"clone"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	new_name := getParam(r, "new_name")
	ip := getParam(r, "ip")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if new_name == "" {
		http.Error(w, "[ERROR]: Missing new_name parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, new_name)
	if ip == "" {
		http.Error(w, "[ERROR]: Missing ip parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, ip)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleCmdHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"cmd"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	command := getParam(r, "command")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if command == "" {
		http.Error(w, "[ERROR]: Missing command parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(command)...)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleConfigHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"config"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	property := getParam(r, "property")
	value := getParam(r, "value")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action != "set" && action != "add" && action != "get" && action != "remove" {
		http.Error(w, "[ERROR]: Unknown action parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, action)
	if property == "" {
		http.Error(w, "[ERROR]: Missing property parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, property)
	if value != "" {
		cmdArgs = append(cmdArgs, value)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleConsoleHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"console"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	user := getParam(r, "user")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if user != "" {
		cmdArgs = append(cmdArgs, user)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleConvertHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"convert"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	release := getParam(r, "release")

	if options != "" {
		options = options + " -ay"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	} else {
		options = "-ay"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if release != "" {
		cmdArgs = append(cmdArgs, release)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleCpHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"cp"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	host_path := getParam(r, "host_path")
	jail_path := getParam(r, "jail_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if host_path == "" {
		http.Error(w, "[ERROR]: Missing host_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, host_path)
	if jail_path == "" {
		http.Error(w, "[ERROR]: Missing jail_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, jail_path)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleCreateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"create"}

	options := getParam(r, "options")
	name := getParam(r, "name")
	release := getParam(r, "release")
	ip := getParam(r, "ip")
	iface := getParam(r, "iface")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if name == "" {
		http.Error(w, "[ERROR]: Missing name paramerter.", http.StatusBadRequest)
		return
	}
	if release == "" {
		http.Error(w, "[ERROR]: Missing release parameter", http.StatusBadRequest)
		return
	}
	if ip == "" {
		http.Error(w, "[ERROR]: Missing ip parameter", http.StatusBadRequest)
		return
	}

	cmdArgs = append(cmdArgs, name, release, ip)

	if iface != "" {
		cmdArgs = append(cmdArgs, iface)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleDestroyHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"destroy"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		options = options + " -ay"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	} else {
		options = "-ay"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleEditHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"edit"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	file := getParam(r, "file")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if file != "" {
		cmdArgs = append(cmdArgs, file)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleEtcupdateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"etcupdate"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	release := getParam(r, "release")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if action == "bootstrap" {
		if release == "" {
			http.Error(w, "[ERROR]: Missing release parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, action, release)
	} else {
		if target == "" {
			http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, target)
		if action == "update" {
			if release == "" {
				http.Error(w, "[ERROR]: Missing release parameter", http.StatusBadRequest)
				return
			}
			cmdArgs = append(cmdArgs, release)
		} else {
			if action == "" {
				http.Error(w, "[ERROR]: Missing action parameter", http.StatusBadRequest)
				return
			}
			cmdArgs = append(cmdArgs, action)
		}
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleExportHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"export"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	path := getParam(r, "path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if path != "" {
		cmdArgs = append(cmdArgs, path)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleHtopHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"htop"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleImportHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"import"}

	options := getParam(r, "options")
	file := getParam(r, "file")
	release := getParam(r, "release")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if file == "" {
		http.Error(w, "[ERROR]: Missing file parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, file)
	if release != "" {
		cmdArgs = append(cmdArgs, release)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleJcpHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"jcp"}

	options := getParam(r, "options")
	source_jail := getParam(r, "source_jail")
	source_path := getParam(r, "source_path")
	destination_jail := getParam(r, "destination_jail")
	destination_path := getParam(r, "destination_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if source_jail == "" {
		http.Error(w, "[ERROR]: Missing source_jail parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, source_jail)
	if source_path == "" {
		http.Error(w, "[ERROR]: Missing source_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, source_path)
	if destination_jail == "" {
		http.Error(w, "[ERROR]: Missing destination_jail parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, destination_jail)
	if destination_path == "" {
		http.Error(w, "[ERROR]: Missing destination_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, destination_path)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleLimitsHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"limits"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	args := getParam(r, "args")
	option := getParam(r, "option")
	value := getParam(r, "value")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "" {
		http.Error(w, "[ERROR]: Missing action parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, action)
	if action == "add" {
		if option == "" {
			http.Error(w, "[ERROR]: Missing option parameter", http.StatusBadRequest)
			return
		}
		if value == "" {
			http.Error(w, "[ERROR]: Missing value parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, option, value)
	} else if action == "remove" {
		if option == "" {
			http.Error(w, "[ERROR]: Missing option parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, option)
	} else if action == "clear" || action == "reset" || action == "stats" {
		cmdArgs = append(cmdArgs, action)
	} else if action == "list" || action == "show" {
		if args == "active" {
			cmdArgs = append(cmdArgs, action, args)
		} else {
			cmdArgs = append(cmdArgs, action)
		}
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleListHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"list"}

	options := getParam(r, "options")
	item := getParam(r, "item")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if item != "" {
		cmdArgs = append(cmdArgs, item)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleMigrateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"migrate"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	destination := getParam(r, "destination")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if destination == "" {
		http.Error(w, "[ERROR]: Missing destination parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, destination)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleMonitorHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"monitor"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	service := getParam(r, "service")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if action == "enable" || action == "disable" || action == "status" {
		cmdArgs = append(cmdArgs, action)
	} else if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "add" || action == "delete" {
		if service == "" {
			http.Error(w, "[ERROR]: Missing service parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, action, service)
	} else if action == "list" {
		cmdArgs = append(cmdArgs, action)
		if service != "" {
			cmdArgs = append(cmdArgs, service)
		}
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleMountHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"mount"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	host_path := getParam(r, "host_path")
	jail_path := getParam(r, "jail_path")
	fs_type := getParam(r, "fs_type")
	fs_options := getParam(r, "fs_options")
	dump := getParam(r, "dump")
	pass_number := getParam(r, "pass_number")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if host_path == "" {
		http.Error(w, "[ERROR]: Missing host_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, host_path)
	if jail_path == "" {
		http.Error(w, "[ERROR]: Missing jail_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, jail_path)
	if fs_type != "" || fs_options != "" || dump != "" || pass_number != "" {
		if fs_type == "" {
			http.Error(w, "[ERROR]: Missing fs_type parameter", http.StatusBadRequest)
			return
		}
		if fs_options == "" {
			http.Error(w, "[ERROR]: Missing fs_options parameter", http.StatusBadRequest)
			return
		}
		if dump == "" {
			http.Error(w, "[ERROR]: Missing dump parameter", http.StatusBadRequest)
			return
		}
		if pass_number == "" {
			http.Error(w, "[ERROR]: Missing pass_number parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, fs_type, fs_options, dump, pass_number)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleNetworkHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"network"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	iface := getParam(r, "iface")
	ip := getParam(r, "ip")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "add" {
		if iface == "" {
			http.Error(w, "[ERROR]: Missing iface parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, action, iface)
		if ip != "" {
			cmdArgs = append(cmdArgs, ip)
		}
	} else {
		if action == "remove" {
			if iface == "" {
				http.Error(w, "[ERROR]: Missing iface parameter", http.StatusBadRequest)
				return
			}
			cmdArgs = append(cmdArgs, action, iface)
		}
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastillePkgHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"pkg"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	args := getParam(r, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleRcpHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"rcp"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	jail_path := getParam(r, "jail_path")
	host_path := getParam(r, "host_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if jail_path == "" {
		http.Error(w, "[ERROR]: Missing jail_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, jail_path)
	if host_path == "" {
		http.Error(w, "[ERROR]: Missing host_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, host_path)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleRdrHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"rdr"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	protocol := getParam(r, "protocol")
	host_port := getParam(r, "host_port")
	jail_port := getParam(r, "jail_port")
	log_options := getParam(r, "log_options")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "clear" || action == "reset" || action == "list" {
		cmdArgs = append(cmdArgs, action)
	} else {
		if protocol == "" {
			http.Error(w, "[ERROR]: Missing protocol parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, protocol)
		if host_port == "" {
			http.Error(w, "[ERROR]: Missing host_port parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, host_port)
		if jail_port == "" {
			http.Error(w, "[ERROR]: Missing jail_port parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, jail_port)
		if action == "log" {
			if log_options == "" {
				http.Error(w, "[ERROR]: Missing log_options parameter", http.StatusBadRequest)
				return
			}
			cmdArgs = append(cmdArgs, action, log_options)
		}
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleRenameHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"rename"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	new_name := getParam(r, "new_name")

	if options != "" {
		options = options + " -a"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	} else {
		options = "-a"
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	if new_name == "" {
		http.Error(w, "[ERROR]: Missing new_name parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target, new_name)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleRestartHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"restart"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleServiceHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"service"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	service := getParam(r, "service")
	args := getParam(r, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if service == "" {
		http.Error(w, "[ERROR]: Missing service parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, service)
	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleSetupHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"setup"}

	options := getParam(r, "options")
	item := getParam(r, "item")
	args := getParam(r, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if item != "" {
		cmdArgs = append(cmdArgs, item)
	}
	if args != "" {
		cmdArgs = append(cmdArgs, args)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleStartHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"start"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleStopHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"stop"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleSysrcHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"sysrc"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	args := getParam(r, "args")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if args == "" {
		http.Error(w, "[ERROR]: Missing args parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, strings.Fields(args)...)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleTagsHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"tags"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	tags := getParam(r, "tags")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "add" || action == "delete" {
		if tags == "" {
			http.Error(w, "[ERROR]: Missing tags parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, action, tags)
	} else if action == "list" {
		cmdArgs = append(cmdArgs, action)
		if tags != "" {
			cmdArgs = append(cmdArgs, tags)
		}
	} else {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleTemplateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"template"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	template := getParam(r, "template")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if action == "convert" {
		cmdArgs = append(cmdArgs, action)
		if template == "" {
			http.Error(w, "[ERROR]: Missing template parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, template)
	} else {
		if target == "" {
			http.Error(w, "[ERROR]: Missing tags parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, target)
		if template == "" {
			http.Error(w, "[ERROR]: Missing template parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, template)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleTopHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"top"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleUmountHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"umount"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	jail_path := getParam(r, "jail_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if jail_path == "" {
		http.Error(w, "[ERROR]: Missing jail_path parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, jail_path)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleUpdateHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"update"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleUpgradeHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"upgrade"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	new_release := getParam(r, "new_release")
	action := getParam(r, "action")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "install" {
		cmdArgs = append(cmdArgs, action)
	} else {
		if new_release == "" {
			http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, new_release)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleVerifyHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"verify"}

	options := getParam(r, "options")
	target := getParam(r, "target")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)

	ParseAndRunCommand(w, r, cmdArgs)
}

func BastilleZfsHandler(w http.ResponseWriter, r *http.Request) {

	cmdArgs := []string{"zfs"}

	options := getParam(r, "options")
	target := getParam(r, "target")
	action := getParam(r, "action")
	tag := getParam(r, "tag")
	key_value := getParam(r, "key_value")
	dataset := getParam(r, "dataset")
	jail_path := getParam(r, "jail_path")

	if options != "" {
		cmdArgs = append(cmdArgs, strings.Fields(options)...)
	}
	if target == "" {
		http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
		return
	}
	cmdArgs = append(cmdArgs, target)
	if action == "snapshot" || action == "destroy" || action == "rollback" {
		cmdArgs = append(cmdArgs, action)
		if tag != "" {
			cmdArgs = append(cmdArgs, tag)
		}
	} else if action == "df" || action == "usage" {
		cmdArgs = append(cmdArgs, action)
	} else if action == "get" || action == "set" {
		cmdArgs = append(cmdArgs, action)
		if key_value == "" {
			http.Error(w, "[ERROR]: Missing target parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, key_value)
	} else if action == "jail" {
		cmdArgs = append(cmdArgs, action)
		if dataset == "" {
			http.Error(w, "[ERROR]: Missing dataset parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, dataset)
		if jail_path == "" {
			http.Error(w, "[ERROR]: Missing jail_path parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, jail_path)
	} else if action == "unjail" {
		cmdArgs = append(cmdArgs, action)
		if jail_path == "" {
			http.Error(w, "[ERROR]: Missing jail_path parameter", http.StatusBadRequest)
			return
		}
		cmdArgs = append(cmdArgs, jail_path)
	}

	ParseAndRunCommand(w, r, cmdArgs)
}
