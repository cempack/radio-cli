package keymap

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up            key.Binding
	Down          key.Binding
	Top           key.Binding
	Bottom        key.Binding
	Search        key.Binding
	Tune          key.Binding
	Favorite      key.Binding
	Related       key.Binding
	Detail        key.Binding
	CopyURL       key.Binding
	AddCrate      key.Binding
	Preset        key.Binding
	PlayPause     key.Binding
	Stop          key.Binding
	Mute          key.Binding
	VolumeUp      key.Binding
	VolumeDown    key.Binding
	CyclePanes    key.Binding
	PrevHist      key.Binding
	NextHist      key.Binding
	Back          key.Binding
	Help          key.Binding
	ViewHome      key.Binding
	ViewDiscover  key.Binding
	ViewFavorites key.Binding
	ViewHistory   key.Binding
	ViewPresets   key.Binding
	ViewCrates    key.Binding
}

var Default = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "bottom"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Tune: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "tune"),
	),
	Favorite: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "favorite"),
	),
	Related: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "related"),
	),
	Detail: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "detail"),
	),
	CopyURL: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy url"),
	),
	AddCrate: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "crate"),
	),
	Preset: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preset"),
	),
	PlayPause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "play/pause"),
	),
	Stop: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "stop"),
	),
	Mute: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "mute"),
	),
	VolumeUp: key.NewBinding(
		key.WithKeys("="),
		key.WithHelp("=", "vol+"),
	),
	VolumeDown: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "vol-"),
	),
	CyclePanes: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "cycle panes"),
	),
	PrevHist: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "prev"),
	),
	NextHist: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next"),
	),
	Back: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "back/quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	ViewHome: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "home"),
	),
	ViewDiscover: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "discover"),
	),
	ViewFavorites: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "favorites"),
	),
	ViewHistory: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "history"),
	),
	ViewPresets: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("5", "presets"),
	),
	ViewCrates: key.NewBinding(
		key.WithKeys("6"),
		key.WithHelp("6", "crates"),
	),
}
