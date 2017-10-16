package execworker

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"syscall"
)

// StartNginx starts secondary nginx process for specific stream
func StartNginx(configPath string, pidPath string) error {
	err := exec.Command("nginx", "-c", configPath, "-g", "'pid "+pidPath+";'").Run()
	if err != nil {
		return err
	}

	return nil
}

// StopNginx accepts a pidPath and kills the process
func StopNginx(pidPath string) error {
	out, err := ioutil.ReadFile(pidPath)
	if err != nil {
		return err
	}

	pid, _ := strconv.Atoi(string(out))

	syscall.Kill(pid, syscall.SIGQUIT)

	return nil
}
