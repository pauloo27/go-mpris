package main

import (
	"log"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
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

	ch := make(chan *dbus.Signal)
	err := player.OnSignal(ch)
	if err != nil {
		panic(err)
	}

	sig := <-ch
	fmt.Println(sig.Body)
}
