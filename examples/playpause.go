package main

import (
	"log"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	names, err := mpris.List(conn)
	if err != nil {
		panic(err)
	}

	if len(names) == 0 {
		log.Fatal("No player found")
	}

	name := names[0]

	player := mpris.New(conn, name)

	status, err := player.GetPlaybackStatus()
	if err != nil {
		log.Fatal("Could not get current playback status")
	}

	log.Printf("The player was %s...", status)

	err = player.PlayPause()
	if err != nil {
		log.Fatal("Could not play/pause player")
	}
}
