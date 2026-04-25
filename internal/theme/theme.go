package theme

import "github.com/charmbracelet/lipgloss"

const (
	ColorBg       = "#1a1714"
	ColorSurface  = "#2d2926"
	ColorText     = "#e8e0d5"
	ColorMuted    = "#8a7f75"
	ColorAccent   = "#d4883a"
	ColorBorder   = "#3d3835"
	ColorSelected = "#1e3a5f"
	ColorSuccess  = "#4a8c5c"
	ColorWarning  = "#d4883a"
	ColorError    = "#8c4a4a"
)

var (
	AppStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorBg)).
		Foreground(lipgloss.Color(ColorText))

	PaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorBorder)).
		Background(lipgloss.Color(ColorBg)).
		Foreground(lipgloss.Color(ColorText))

	FocusedPaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorAccent)).
		Background(lipgloss.Color(ColorBg)).
		Foreground(lipgloss.Color(ColorText))

	TitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorAccent)).
		Bold(true).
		Padding(0, 1)

	SubtleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted))

	BadgeStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorMuted)).
		Padding(0, 1)

	AccentBadgeStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorAccent)).
		Foreground(lipgloss.Color(ColorBg)).
		Padding(0, 1)

	SelectedRowStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSelected)).
		Foreground(lipgloss.Color(ColorText)).
		Bold(true)

	NowPlayingTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorAccent)).
		Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSurface)).
		Foreground(lipgloss.Color(ColorMuted)).
		Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorError)).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess))
)
