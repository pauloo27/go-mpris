package mpris

import (
	"fmt"
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

func getProperty(obj *dbus.Object, iface string, prop string) (dbus.Variant, error) {
	result := dbus.Variant{}
	err := obj.Call(getPropertyMethod, 0, iface, prop).Store(&result)
	if err != nil {
		return dbus.Variant{}, err
	}
	return result, nil
}

func setProperty(obj *dbus.Object, iface string, prop string, val interface{}) error {
	call := obj.Call(setPropertyMethod, 0, iface, prop, dbus.MakeVariant(val))
	return call.Err
}

func convertToMicroseconds(seconds float64) int64 {
	return int64(seconds * 1000000)
}

func convertToSeconds(microseconds int64) float64 {
	return float64(microseconds) / 1000000.0
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
func (i *base) Raise() error {
	return i.obj.Call(BaseInterface+".Raise", 0).Err
}

// Quit closes the player.
func (i *base) Quit() error {
	return i.obj.Call(BaseInterface+".Quit", 0).Err
}

// GetIdentity returns the player identity.
func (i *base) GetIdentity() (string, error) {
	value, err := getProperty(i.obj, BaseInterface, "Identity")

	return value.Value().(string), err
}

type player struct {
	obj *dbus.Object
}

// Next skips to the next track in the tracklist.
func (i *player) Next() error {
	return i.obj.Call(PlayerInterface+".Next", 0).Err
}

// Previous skips to the previous track in the tracklist.
func (i *player) Previous() error {
	return i.obj.Call(PlayerInterface+".Previous", 0).Err
}

// Pause pauses the current track.
func (i *player) Pause() error {
	return i.obj.Call(PlayerInterface+".Pause", 0).Err
}

// PlayPause resumes the current track if it's paused and pauses it if it's playing.
func (i *player) PlayPause() error {
	return i.obj.Call(PlayerInterface+".PlayPause", 0).Err
}

// Stop stops the current track.
func (i *player) Stop() error {
	return i.obj.Call(PlayerInterface+".Stop", 0).Err
}

// Play starts or resumes the current track.
func (i *player) Play() error {
	return i.obj.Call(PlayerInterface+".Play", 0).Err
}

// Seek seeks the current track position by the offset. The offset should be in seconds.
// If the offset is negative it's seeked back.
func (i *player) Seek(offset float64) error {
	return i.obj.Call(PlayerInterface+".Seek", 0, convertToMicroseconds(offset)).Err
}

// SetTrackPosition sets the position of a track. The position should be in seconds.
func (i *player) SetTrackPosition(trackId *dbus.ObjectPath, position float64) error {
	return i.obj.Call(PlayerInterface+".SetPosition", 0, trackId, convertToMicroseconds(position)).Err
}

// OpenUri opens and plays the uri if supported.
func (i *player) OpenUri(uri string) error {
	return i.obj.Call(PlayerInterface+".OpenUri", 0, uri).Err
}

// PlaybackStatus the status of the playback. It can be "Playing", "Paused" or "Stopped".
type PlaybackStatus string

const (
	PlaybackPlaying PlaybackStatus = "Playing"
	PlaybackPaused  PlaybackStatus = "Paused"
	PlaybackStopped PlaybackStatus = "Stopped"
)

// GetPlaybackStatus gets the playback status.
func (i *player) GetPlaybackStatus() (PlaybackStatus, error) {
	variant, err := i.obj.GetProperty(PlayerInterface + ".PlaybackStatus")
	if err != nil {
		return "", err
	}
	if variant.Value() == nil {
		return "", fmt.Errorf("Variant value is nil")
	}
	return PlaybackStatus(variant.Value().(string)), nil
}

// LoopStatus the status of the player loop. It can be "None", "Track" or "Playlist".
type LoopStatus string

const (
	LoopNone     LoopStatus = "None"
	LoopTrack    LoopStatus = "Track"
	LoopPlaylist LoopStatus = "Playlist"
)

// GetLoopStatus returns the loop status.
func (i *player) GetLoopStatus() (LoopStatus, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "LoopStatus")
	if err != nil {
		return LoopStatus(""), err
	}
	if variant.Value() == nil {
		return "", fmt.Errorf("Variant value is nil")
	}
	return LoopStatus(variant.Value().(string)), nil
}

// SetLoopStatus sets the loop status to loopStatus.
func (i *player) SetLoopStatus(loopStatus LoopStatus) error {
	return i.SetPlayerProperty("LoopStatus", loopStatus)
}

// SetProperty sets the value of a propertyName in the targetInterface.
func (i *player) SetProperty(targetInterface, propertyName string, value interface{}) error {
	return setProperty(i.obj, targetInterface, propertyName, value)
}

// SetPlayerProperty sets the propertyName from the player interface.
func (i *player) SetPlayerProperty(propertyName string, value interface{}) error {
	return setProperty(i.obj, PlayerInterface, propertyName, value)
}

// GetProperty returns the properityName in the targetInterface.
func (i *player) GetProperty(targetInterface, properityName string) (dbus.Variant, error) {
	return getProperty(i.obj, targetInterface, properityName)
}

// GetPlayerProperty returns the properityName from the player interface.
func (i *player) GetPlayerProperty(properityName string) (dbus.Variant, error) {
	return getProperty(i.obj, PlayerInterface, properityName)
}

// Returns the current playback rate.
func (i *player) GetRate() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Rate")
	if err != nil {
		return 0.0, err
	}
	if variant.Value() == nil {
		return 0.0, fmt.Errorf("Variant value is nil")
	}
	return variant.Value().(float64), nil
}

// GetShuffle returns false if the player is going linearly through a playlist and false if it's
// in some other order.
func (i *player) GetShuffle() (bool, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Shuffle")
	if err != nil {
		return false, err
	}
	if variant.Value() == nil {
		return false, fmt.Errorf("Variant value is nil")
	}
	return variant.Value().(bool), nil
}

// GetMetadata returns the metadata.
func (i *player) GetMetadata() (map[string]dbus.Variant, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Metadata")
	if err != nil {
		return nil, err
	}
	if variant.Value() == nil {
		return nil, fmt.Errorf("Variant value is nil")
	}
	return variant.Value().(map[string]dbus.Variant), nil
}

// GetVolume returns the volume.
func (i *player) GetVolume() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Volume")
	if err != nil {
		return 0.0, err
	}
	if variant.Value() == nil {
		return 0.0, fmt.Errorf("Variant value is nil")
	}
	return variant.Value().(float64), nil
}

// SetVolume sets the volume.
func (i *player) SetVolume(volume float64) error {
	return setProperty(i.obj, PlayerInterface, "Volume", volume)
}

// GetLength returns the current track length in seconds.
func (i *player) GetLength() (float64, error) {
	metadata, err := i.GetMetadata()
	if err != nil {
		return 0.0, err
	}
	if metadata == nil || metadata["mpris:length"].Value() == nil {
		return 0.0, fmt.Errorf("Variant value is nil")
	}
	return convertToSeconds(metadata["mpris:length"].Value().(int64)), nil
}

// GetPosition returns the position in seconds of the current track.
func (i *player) GetPosition() (float64, error) {
	variant, err := getProperty(i.obj, PlayerInterface, "Position")
	if err != nil {
		return 0.0, err
	}
	if variant.Value() == nil {
		return 0.0, fmt.Errorf("Variant value is nil")
	}
	return convertToSeconds(variant.Value().(int64)), nil
}

// SetPosition sets the position of the current track. The position should be in seconds.
func (i *player) SetPosition(position float64) error {
	metadata, err := i.GetMetadata()
	if err != nil {
		return err
	}
	if metadata == nil || metadata["mpris:trackid"].Value() == nil {
		return fmt.Errorf("Variant value is nil")
	}
	trackId := metadata["mpris:trackid"].Value().(dbus.ObjectPath)
	i.SetTrackPosition(&trackId, position)
	return nil
}

// New connects the the player with the name in the connection conn.
func New(conn *dbus.Conn, name string) *Player {
	obj := conn.Object(name, dbusObjectPath).(*dbus.Object)

	return &Player{
		&base{obj},
		&player{obj},
	}
}
