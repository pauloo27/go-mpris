# GO-MPRIS

A Go library for MPRIS.

## Install

> $ go get github.com/Pauloo27/go-mpris

_the dependency github.com/godbus/dbus/v5 is going to be installed as well._


## Example
Printing the current playback status and then changing it:
```go
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
```

**For more examples, see the [examples folder](./examples).**

## Go Docs
Read the docs at https://pkg.go.dev/github.com/Pauloo27/go-mpris.

## Credits
[emersion](https://github.com/emersion/go-mpris) for the original code.

