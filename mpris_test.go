package mpris

import (
	"testing"
	"time"

	"github.com/godbus/dbus"
)

func checkVolume(t *testing.T, player *Player) {
	volume := player.GetVolume()
	t.Logf("Current player volume %f", volume)
	player.SetVolume(0.5)
	time.Sleep(1 * time.Second)
	player.SetVolume(volume)
}

func checkPlayback(t *testing.T, player *Player) {
	status := player.GetPlaybackStatus()

	if status != PlaybackPlaying && status != PlaybackStopped && status != PlaybackPaused {
		t.Errorf("%s is not a valid playback status", status)
	} else {
		t.Logf("Player playback status is %s", status)
	}

}

func checkLoop(t *testing.T, player *Player) {
	if !player.HasLoopStatus() {
		t.Logf("Player don't have a loop status")
		return
	}
	loopStatus := player.GetLoopStatus()

	if loopStatus != LoopNone && loopStatus != LoopTrack && loopStatus != LoopPlaylist {
		t.Errorf("%s is not a valid loop status", loopStatus)
	} else {
		t.Logf("Players loop status is %s", loopStatus)
	}
}

func TestPlayer(t *testing.T) {
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

	t.Run("Playback", func(t *testing.T) { checkPlayback(t, player) })
	t.Run("Loop", func(t *testing.T) { checkLoop(t, player) })
	t.Run("Volume", func(t *testing.T) { checkVolume(t, player) })
}
