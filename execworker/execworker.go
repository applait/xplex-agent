package execworker

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strconv"
	"syscall"
)

func StartStreamer(streamer string, configPath string, pidPath string) error {
	switch streamer {
	case "nginx":
		err := startNginx(configPath, pidPath)
		if err != nil {
			return nil
		}
	default:
		return errors.New("Unknown streamer")
	}

	return nil
}

func StopStreamer(streamer string, pidPath string) error {
	switch streamer {
	case "nginx":
		err := stopNginx(pidPath)
		if err != nil {
			return nil
		}
	default:
		return errors.New("Unknown streamer")
	}

	return nil
}

// startNginx starts secondary nginx process for specific stream
func startNginx(configPath string, pidPath string) error {
	err := exec.Command("/usr/local/nginx/sbin/nginx", "-c", configPath, "-g", "'pid "+pidPath+";'").Run()
	if err != nil {
		return err
	}

	return nil
}

// stopNginx accepts a pidPath and kills the process
func stopNginx(pidPath string) error {
	out, err := ioutil.ReadFile(pidPath)
	if err != nil {
		return err
	}

	pid, _ := strconv.Atoi(string(out))

	syscall.Kill(pid, syscall.SIGQUIT)

	return nil
}
