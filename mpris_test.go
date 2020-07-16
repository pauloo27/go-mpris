package mpris

import (
	"testing"

	"github.com/godbus/dbus"
)

func TestPlaybackStatus(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Error(err)
	}

	names, err := List(conn)
	if err != nil {
		t.Error(err)
	}
	if len(names) == 0 {
		t.Error("No players found")
		return
	}

	name := names[0]
	t.Logf("Found player %s", name)

	player := New(conn, name)

	status := player.GetPlaybackStatus()

	if status != PlaybackPlaying && status != PlaybackStopped && status != PlaybackPaused {
		t.Errorf("%s is not a valid playback status", status)
	} else {
		t.Logf("Player %s playback status is %s", name, player.GetPlaybackStatus())
	}
}
