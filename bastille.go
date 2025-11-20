package main

import (
	"fmt"
	"os/exec"
)

func runBastille(args ...string) error {
	cmd := exec.Command("bastille", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bastille %v failed: %v\n%s", args, err, output)
	}
	return nil
}

func BastilleCreate(name, release, ip, iface string) error {
	args := []string{"create", name, release, ip}
	if iface != "" {
		args = append(args, iface)
	}
	return runBastille(args...)
}

func BastilleDestroy(name string) error  { return runBastille("destroy","-y", name) }
func BastilleStart(name string) error    { return runBastille("start", name) }
func BastilleStop(name string) error     { return runBastille("stop", name) }
func BastilleRestart(name string) error  { return runBastille("restart", name) }
func BastilleRename(old, new string) error { return runBastille("rename", old, new) }
func BastilleUpgrade(name string) error  { return runBastille("upgrade", name) }
func BastilleList() (string, error)  {
	cmd := exec.Command("bastille", "list", "all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("bastille list failed: %v\n%s", err, output)
	}
	return string(output), nil
}

