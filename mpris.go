package mpris

import (
	"log"
	"strings"

	"github.com/godbus/dbus"
)

const (
	dbusObjectPath          = "/org/mpris/MediaPlayer2"
	propertiesChangedSignal = "org.freedesktop.DBus.Properties.PropertiesChanged"

	BaseInterface      = "org.mpris.MediaPlayer2"
	PlayerInterface    = "org.mpris.MediaPlayer2.Player"
	TrackListInterface = "org.mpris.MediaPlayer2.TrackList"
	PlaylistsInterface = "org.mpris.MediaPlayer2.Playlists"

	getPropertyMethod = "org.freedesktop.DBus.Properties.Get"
	setPropertyMethod = "org.freedesktop.DBus.Properties.Set"
)

func getProperty(obj *dbus.Object, iface string, prop string) dbus.Variant {
	result := dbus.Variant{}
	err := obj.Call(getPropertyMethod, 0, iface, prop).Store(&result)
	if err != nil {
		log.Println("Warning: could not get dbus property:", iface, prop, err)
		return dbus.Variant{}
	}
	return result
}

func setProperty(obj *dbus.Object, iface string, prop string, val interface{}) {
	call := obj.Call(setPropertyMethod, 0, prop, val)
	if call.Err != nil {
		log.Println("Warning: could not set dbus property:", iface, prop, call.Err)
	}
}

func convertToMicroseconds(seconds float64) int64 {
	return int64(seconds * 1000000)
}

// List lists the available players.
func List(conn *dbus.Conn) ([]string, error) {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return nil, err
	}

	var mprisNames []string
	for _, name := range names {
		if strings.HasPrefix(name, BaseInterface) {
			mprisNames = append(mprisNames, name)
		}
	}
	return mprisNames, nil
}

// Player represents a mpris player.
type Player struct {
	*base
	*player
}

type base struct {
	obj *dbus.Object
}

// Raise raises player priority.
func (i *base) Raise() {
	i.obj.Call(BaseInterface+".Raise", 0)
}

// Quit closes the player.
func (i *base) Quit() {
	i.obj.Call(BaseInterface+".Quit", 0)
}

// GetIdentity returns the player identity.
func (i *base) GetIdentity() string {
	return getProperty(i.obj, BaseInterface, "Identity").Value().(string)
}

type player struct {
	obj *dbus.Object
}

// Next skips to the next track in the tracklist.
func (i *player) Next() {
	i.obj.Call(PlayerInterface+".Next", 0)
}

// Previous skips to the previous track in the tracklist.
func (i *player) Previous() {
	i.obj.Call(PlayerInterface+".Previous", 0)
}

// Pause pauses the current track.
func (i *player) Pause() {
	i.obj.Call(PlayerInterface+".Pause", 0)
}

// PlayPause resumes the current track if it's paused and pauses it if it's playing.
func (i *player) PlayPause() {
	i.obj.Call(PlayerInterface+".PlayPause", 0)
}

// Stop stops the current track.
func (i *player) Stop() {
	i.obj.Call(PlayerInterface+".Stop", 0)
}

// Play starts or resumes the current track.
func (i *player) Play() {
	i.obj.Call(PlayerInterface+".Play", 0)
}

// Seek seeks the current track position by the offset. The offset should be in seconds.
// If the offset is negative it's seeked back.
func (i *player) Seek(offset float64) {
	i.obj.Call(PlayerInterface+".Seek", 0, convertToMicroseconds(offset))
}

// SetTrackPosition sets the position of a track.
func (i *player) SetTrackPosition(trackId *dbus.ObjectPath, position float64) {
	i.obj.Call(PlayerInterface+".SetPosition", 0, trackId, convertToMicroseconds(position))
}

// OpenUri opens and plays the uri if supported.
func (i *player) OpenUri(uri string) {
	i.obj.Call(PlayerInterface+".OpenUri", 0, uri)
}

// PlaybackStatus the status of the playback. It can be "Playing", "Paused" or "Stopped".
type PlaybackStatus string

const (
	PlaybackPlaying PlaybackStatus = "Playing"
	PlaybackPaused  PlaybackStatus = "Paused"
	PlaybackStopped PlaybackStatus = "Stopped"
)

// GetPlaybackStatus gets the playback status.
func (i *player) GetPlaybackStatus() PlaybackStatus {
	variant, err := i.obj.GetProperty(PlayerInterface + ".PlaybackStatus")
	if err != nil {
		return ""
	}
	return PlaybackStatus(variant.Value().(string))
}

// LoopStatus the status of the player loop. It can be "None", "Track" or "Playlist".
type LoopStatus string

const (
	LoopNone     LoopStatus = "None"
	LoopTrack    LoopStatus = "Track"
	LoopPlaylist LoopStatus = "Playlist"
)

// HasLookStatus returns if the player support loop status.
func (i *player) HasLoopStatus() bool {
	return getProperty(i.obj, PlayerInterface, "LoopStatus").Value() != nil
}

// GetLoopStatus returns the loop status.
func (i *player) GetLoopStatus() LoopStatus {
	return LoopStatus(getProperty(i.obj, PlayerInterface, "LoopStatus").Value().(string))
}

// GetProperty returns the properityName in the targetInterface.
func (i *player) GetProperty(targetInterface, properityName string) dbus.Variant {
	return getProperty(i.obj, targetInterface, properityName)
}

// GetPlayerProperty returns the properityName from the player interface.
func (i *player) GetPlayerProperty(properityName string) dbus.Variant {
	return getProperty(i.obj, PlayerInterface, properityName)
}

// Returns the current playback rate.
func (i *player) GetRate() float64 {
	return getProperty(i.obj, PlayerInterface, "Rate").Value().(float64)
}

// GetShuffle returns false if the player is going linearly through a playlist and false if it's
// in some other order.
func (i *player) GetShuffle() bool {
	return getProperty(i.obj, PlayerInterface, "Shuffle").Value().(bool)
}

// GetMetadata returns the metadata.
func (i *player) GetMetadata() map[string]dbus.Variant {
	return getProperty(i.obj, PlayerInterface, "Metadata").Value().(map[string]dbus.Variant)
}

// GetVolume returns the volume.
func (i *player) GetVolume() float64 {
	return getProperty(i.obj, PlayerInterface, "Volume").Value().(float64)
}

// SetVolume sets the volume.
func (i *player) SetVolume(volume float64) {
	setProperty(i.obj, PlayerInterface, "Volume", volume)
}

// GetPosition returns the position of the current track.
func (i *player) GetPosition() int64 {
	return getProperty(i.obj, PlayerInterface, "Position").Value().(int64)
}

// SetPosition sets the position of the current track.
func (i *player) SetPosition(position float64) {
	trackId := i.GetMetadata()["mpris:trackid"].Value().(dbus.ObjectPath)
	i.SetTrackPosition(&trackId, position)
}

// New connects the the player with the name in the connection conn.
func New(conn *dbus.Conn, name string) *Player {
	obj := conn.Object(name, dbusObjectPath).(*dbus.Object)

	return &Player{
		&base{obj},
		&player{obj},
	}
}
