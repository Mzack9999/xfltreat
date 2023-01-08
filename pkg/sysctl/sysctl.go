package sysctl

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	sysctlBase = "/proc/sys"
)

// GetSysctl returns the value for the specified sysctl setting
func Get(sysctl string) (int, error) {
	return getWithFile(sysctl)
}

func getWithFile(sysctl string) (int, error) {
	data, err := os.ReadFile(path.Join(sysctlBase, sysctl))
	if err != nil {
		return -1, err
	}
	return toInt(string(data))
}

func toInt(str string) (int, error) {
	val, err := strconv.Atoi(normalize(str))
	if err != nil {
		return -1, err
	}
	return val, nil
}

func getWithExec(sysctl string) (int, error) {
	stdout, stderr, err := run("sysctl", sysctl)
	if err != nil {
		return -1, err
	}
	return stdsToInt(stdout, stderr)
}

// SetSysctl modifies the specified sysctl flag to the new value
func Set(sysctl string, value int) error {
	return setWithFile(sysctl, value)
}

func setWithFile(sysctl string, newValue int) error {
	return os.WriteFile(path.Join(sysctlBase, sysctl), []byte(strconv.Itoa(newValue)), 0640)
}

func setWithExec(sysctl string, newValue int) (int, error) {
	newVar := fmt.Sprintf("%s=%d", path.Join(sysctlBase, sysctl), newValue)
	stdout, stderr, err := run("sysctl", "-w", newVar)
	if err != nil {
		return -1, err
	}
	return stdsToInt(stdout, stderr)
}

func reload() error {
	_, stderr, err := run("sysctl", "-p")
	if len(stderr) > 0 {
		return errors.Wrap(err, stderr)
	}
	return nil
}

func normalize(str string) string {
	return strings.Trim(str, " \n")
}

func run(args ...string) (string, string, error) {
	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	return outb.String(), errb.String(), err
}

func stdsToInt(stdout, stderr string) (int, error) {
	value, err := toInt(stdout)
	if len(stderr) > 0 {
		err = errors.Wrap(err, stderr)
	}
	return value, err
}
