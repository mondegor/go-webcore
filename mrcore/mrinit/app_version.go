package mrinit

import (
	"os/exec"
	"strings"
)

// Version - возвращает версию приложения из внешнего окружения.
func Version() string {
	if _, err := exec.LookPath("git"); err != nil {
		return ""
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	b, err := cmd.Output()
	if err != nil {
		return ""
	}

	value := strings.TrimSpace(string(b))

	// если указана любая ветка кроме мастера
	if value != "master" && value != "main" && value != "HEAD" {
		return value
	}

	// Примеры тегов:
	//   v0.14.7-0-de1493e0
	//   v0.8.1-0-gd3a5efc-dirty
	cmd = exec.Command("git", "describe", "--long", "--always", "--dirty")

	b, err = cmd.Output()
	if err != nil {
		return value
	}

	return strings.TrimSpace(string(b))
}
