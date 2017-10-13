package execworker

import (
	"os/exec"
)

// SpinNginx starts secondary nginx process for specific stream
func SpinNginx(streamConfigPath string) (string, error) {
	out, err := exec.Command("nginx", "-c", streamConfigPath).Output()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
