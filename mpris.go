package mpris

import (
	"log"
	"strings"

	"github.com/godbus/dbus"
)

const (
	dbusObjectPath          = "/org/mpris/MediaPlayer2"
	propertiesChangedSignal = "org.freedesktop.DBus.Properties.PropertiesChanged"

	baseInterface      = "org.mpris.MediaPlayer2"
	playerInterface    = "org.mpris.MediaPlayer2.Player"
	trackListInterface = "org.mpris.MediaPlayer2.TrackList"
	playlistsInterface = "org.mpris.MediaPlayer2.Playlists"

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

func List(conn *dbus.Conn) ([]string, error) {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return nil, err
	}

	var mprisNames []string
	for _, name := range names {
		if strings.HasPrefix(name, baseInterface) {
			mprisNames = append(mprisNames, name)
		}
	}
	return mprisNames, nil
}

type Player struct {
	*base
	*player
}

type base struct {
	obj *dbus.Object
}

func (i *base) Raise() {
	i.obj.Call(baseInterface+".Raise", 0)
}

func (i *base) Quit() {
	i.obj.Call(baseInterface+".Quit", 0)
}

func (i *base) GetIdentity() string {
	return getProperty(i.obj, baseInterface, "Identity").String()
}

type player struct {
	obj *dbus.Object
}

func (i *player) Next() {
	i.obj.Call(playerInterface+".Next", 0)
}

func (i *player) Previous() {
	i.obj.Call(playerInterface+".Previous", 0)
}

func (i *player) Pause() {
	i.obj.Call(playerInterface+".Pause", 0)
}

func (i *player) PlayPause() {
	i.obj.Call(playerInterface+".PlayPause", 0)
}

func (i *player) Stop() {
	i.obj.Call(playerInterface+".Stop", 0)
}

func (i *player) Play() {
	i.obj.Call(playerInterface+".Play", 0)
}

func (i *player) Seek(offset int64) {
	i.obj.Call(playerInterface+".Seek", 0, offset)
}

func (i *player) SetTrackPosition(trackId *dbus.ObjectPath, position int64) {
	i.obj.Call(playerInterface+".SetPosition", 0, trackId, position)
}

func (i *player) OpenUri(uri string) {
	i.obj.Call(playerInterface+".OpenUri", 0, uri)
}

type PlaybackStatus string

const (
	PlaybackPlaying PlaybackStatus = "Playing"
	PlaybackPaused  PlaybackStatus = "Paused"
	PlaybackStopped PlaybackStatus = "Stopped"
)

func (i *player) GetPlaybackStatus() PlaybackStatus {
	variant, err := i.obj.GetProperty(playerInterface + ".PlaybackStatus")
	if err != nil {
		return ""
	}
	return PlaybackStatus(variant.Value().(string))
}

type LoopStatus string

const (
	LoopNone     LoopStatus = "None"
	LoopTrack    LoopStatus = "Track"
	LoopPlaylist LoopStatus = "Playlist"
)

func (i *player) GetLoopStatus() LoopStatus {
	return LoopStatus(getProperty(i.obj, playerInterface, "LoopStatus").Value().(string))
}

func (i *player) GetRate() float64 {
	return getProperty(i.obj, playerInterface, "Rate").Value().(float64)
}

func (i *player) GetShuffle() bool {
	return getProperty(i.obj, playerInterface, "Shuffle").Value().(bool)
}

func (i *player) GetMetadata() map[string]dbus.Variant {
	return getProperty(i.obj, playerInterface, "Metadata").Value().(map[string]dbus.Variant)
}

func (i *player) GetVolume() float64 {
	return getProperty(i.obj, playerInterface, "Volume").Value().(float64)
}
func (i *player) SetVolume(volume float64) {
	setProperty(i.obj, playerInterface, "Volume", volume)
}

func (i *player) GetPosition() int64 {
	return getProperty(i.obj, playerInterface, "Position").Value().(int64)
}
func (i *player) SetPosition(position float64) {
	setProperty(i.obj, playerInterface, "Position", position)
}

func New(conn *dbus.Conn, name string) *Player {
	obj := conn.Object(name, dbusObjectPath).(*dbus.Object)

	return &Player{
		&base{obj},
		&player{obj},
	}
}
