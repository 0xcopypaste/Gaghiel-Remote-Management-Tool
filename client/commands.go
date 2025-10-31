package main

import (
	"os/exec"
	"strings"
)

func execCommand(cmdStr string) string {
	if strings.TrimSpace(cmdStr) == "" {
		return "No command provided"
	}

	c := exec.Command("cmd.exe", "/C", cmdStr)
	out, err := c.CombinedOutput()
	result := strings.TrimSpace(string(out))

	if err != nil {
		return "Error: " + err.Error()
	}
	if result == "" {
		return "Output: (no output)"
	}
	return "Output:\n" + result
}
