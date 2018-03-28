package execworker

import (
	"errors"
	"os/exec"
)

func StartStreamer(streamer string, streamKey string) error {
	switch streamer {
	case "nginx":
		err := startNginx(streamKey)
		if err != nil {
			return nil
		}
	default:
		return errors.New("Unknown streamer")
	}

	return nil
}

func StopStreamer(streamer string, streamKey string) error {
	switch streamer {
	case "nginx":
		err := stopNginx(streamKey)
		if err != nil {
			return nil
		}
	default:
		return errors.New("Unknown streamer")
	}

	return nil
}

// startNginx starts secondary nginx process for specific stream
func startNginx(streamKey string) error {
	err := exec.Command("systemctl start nginx-" + streamKey).Run()
	if err != nil {
		return err
	}

	return nil
}

// stopNginx accepts a pidPath and kills the process
func stopNginx(streamKey string) error {
	err := exec.Command("systemctl stop nginx-" + streamKey).Run()
	if err != nil {
		return err
	}

	return nil
}
