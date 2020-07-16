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
		log.Fatal("No media player found.")
	}

	name := names[0]
	log.Println("Found media player:", name)

	player := mpris.New(conn, name)

	log.Println("Media player identity:", player.GetIdentity())

	player.Raise()
}
