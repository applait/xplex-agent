package execworker

import (
	"errors"

	"log"

	"github.com/coreos/go-systemd/dbus"
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
	// Init dbus connection
	conn, err := dbus.New()
	if err != nil {
		return err
	}

	// Daemon reload
	err = conn.Reload()
	if err != nil {
		return err
	}

	// Create a channel for reception of nginx success
	c := make(chan string)

	// Start unit
	_, err = conn.StartUnit("nginx-"+streamKey+".service", "replace", c)

	status := <-c
	log.Println("Nginx start status " + streamKey + ": " + status)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// stopNginx accepts a pidPath and kills the process
func stopNginx(streamKey string) error {
	// Init dbus connection
	conn, err := dbus.New()
	if err != nil {
		return err
	}

	// Create a channel for reception of nginx success
	c := make(chan string)

	// Stop unit
	_, err = conn.StopUnit("nginx-"+streamKey+".service", "replace", nil)

	status := <-c
	log.Println("Nginx stop status " + streamKey + ": " + status)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
