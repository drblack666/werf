package utils

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunCommand(dir, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	return cmd.CombinedOutput()
}

func RunSucceedCommand(dir, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter

	errorDesc := fmt.Sprintf("%[2]s %[3]s (dir: %[1]s)", dir, command, strings.Join(args, " "))
	Ω(cmd.Run()).ShouldNot(HaveOccurred(), errorDesc)
}

func SucceedCommandOutput(dir, command string, args ...string) string {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	errorDesc := fmt.Sprintf("%[2]s %[3]s (dir: %[1]s)", dir, command, strings.Join(args, " "))
	res, err := cmd.CombinedOutput()
	Ω(err).ShouldNot(HaveOccurred(), errorDesc)
	return string(res)
}
