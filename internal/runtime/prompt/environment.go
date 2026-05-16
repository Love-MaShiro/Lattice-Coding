package prompt

import (
	"os"
	"runtime"
	"time"
)

type Environment struct {
	WorkingDir string
	OS         string
	Arch       string
	Time       string
	Shell      string
}

func LoadEnvironment(workingDir string) Environment {
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	return Environment{
		WorkingDir: workingDir,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Time:       time.Now().Format(time.RFC3339),
		Shell:      detectShell(),
	}
}

func detectShell() string {
	for _, key := range []string{"SHELL", "ComSpec"} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	return "sh"
}
