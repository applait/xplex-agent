package cron

import (
	"github.com/jasonlvhit/gocron"
)

func pollRig() {
	// TODO: Poll rig, parse JSON and take action
}

// Start starts the CRON which polls rig and take necessary actions on streams
func Start() {
	gocron.Every(2).Minutes().Do(pollRig)

	gocron.Start()
}
