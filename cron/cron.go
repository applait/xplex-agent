package cron

import (
	"github.com/jasonlvhit/gocron"
)

func pollRig() {
	// Fetch existing list of streams, ping rig with the info and wait for action. Ideally a HTTP client
	// On response, get relevant info from boltDB and call execworker
	// TODO: Poll rig, parse JSON and take action
}

// Start starts the CRON which polls rig and take necessary actions on streams
func Start() {
	gocron.Every(5).Minutes().Do(pollRig)

	gocron.Start()
}
