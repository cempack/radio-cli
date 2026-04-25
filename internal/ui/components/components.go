package components

import (
	"fmt"
	"strings"

	"github.com/cempack/radio-cli/internal/domain"
	"github.com/cempack/radio-cli/internal/player"
	"github.com/cempack/radio-cli/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

func StationRow(s domain.Station, selected bool, width int) string {
	health := "●"
	healthColor := theme.ColorSuccess
	if !s.LastCheckOK {
		health = "○"
		healthColor = theme.ColorError
	}
	healthStr := lipgloss.NewStyle().Foreground(lipgloss.Color(healthColor)).Render(health)

	name := s.Name
	maxName := width - 30
	if maxName < 10 {
		maxName = 10
	}
	if len(name) > maxName {
		name = name[:maxName-1] + "…"
	}

	country := ""
	if s.CountryCode != "" {
		country = theme.BadgeStyle.Render(s.CountryCode)
	} else if s.Country != "" && len(s.Country) >= 2 {
		country = theme.BadgeStyle.Render(s.Country[:2])
	}

	bitrate := ""
	if s.Bitrate > 0 {
		bitrate = theme.SubtleStyle.Render(fmt.Sprintf("%dk", s.Bitrate))
	}

	tag := ""
	if len(s.Tags) > 0 {
		tag = theme.SubtleStyle.Render(s.Tags[0])
	}

	row := fmt.Sprintf("%s %-*s %s %s %s", healthStr, maxName, name, country, bitrate, tag)

	if selected {
		return theme.SelectedRowStyle.Width(width).Render(row)
	}
	return lipgloss.NewStyle().Width(width).Render(row)
}

func NowPlayingPanel(s *domain.Station, p *player.Player, width, height int) string {
	if s == nil || p == nil {
		idle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ColorMuted)).
			Width(width).
			Align(lipgloss.Center).
			Render("\n\nNo station playing\n\nPress / to search\nor 2 to discover stations")
		return theme.PaneStyle.Width(width).Height(height).Render(idle)
	}

	var lines []string
	lines = append(lines, theme.TitleStyle.Render("NOW PLAYING"))
	lines = append(lines, "")

	name := s.Name
	if len(name) > width-4 {
		name = name[:width-5] + "…"
	}
	lines = append(lines, theme.NowPlayingTitleStyle.Width(width-4).Render(name))
	lines = append(lines, "")

	stateStr := ""
	stateColor := theme.ColorMuted
	switch p.State() {
	case player.StatePlaying:
		stateStr = "▶ LIVE"
		stateColor = theme.ColorSuccess
	case player.StateConnecting:
		stateStr = "◌ CONNECTING"
		stateColor = theme.ColorWarning
	case player.StateStopped:
		stateStr = "■ STOPPED"
		stateColor = theme.ColorMuted
	case player.StateError:
		stateStr = "✗ ERROR"
		stateColor = theme.ColorError
	}
	stateBadge := lipgloss.NewStyle().
		Background(lipgloss.Color(stateColor)).
		Foreground(lipgloss.Color(theme.ColorBg)).
		Padding(0, 1).
		Render(stateStr)
	lines = append(lines, stateBadge)
	lines = append(lines, "")

	if p.State() == player.StateError && p.Error() != "" {
		lines = append(lines, theme.ErrorStyle.Render(p.Error()))
		lines = append(lines, "")
	}

	if s.Country != "" {
		lines = append(lines, theme.SubtleStyle.Render("Country: ")+s.Country)
	}
	if s.Language != "" {
		lines = append(lines, theme.SubtleStyle.Render("Language: ")+s.Language)
	}
	if s.Codec != "" {
		lines = append(lines, theme.SubtleStyle.Render("Codec: ")+s.Codec)
	}
	if s.Bitrate > 0 {
		lines = append(lines, theme.SubtleStyle.Render("Bitrate: ")+fmt.Sprintf("%d kbps", s.Bitrate))
	}
	if len(s.Tags) > 0 {
		maxTagsToDisplay := 3
		if len(s.Tags) < maxTagsToDisplay {
			maxTagsToDisplay = len(s.Tags)
		}
		lines = append(lines, "")
		tagLine := ""
		for i := 0; i < maxTagsToDisplay; i++ {
			tagLine += theme.BadgeStyle.Render(s.Tags[i]) + " "
		}
		lines = append(lines, strings.TrimSpace(tagLine))
	}
	lines = append(lines, "")
	lines = append(lines, theme.SubtleStyle.Render("x=stop  f=fav  r=related"))

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func StatusBar(width int, statusMsg, errMsg string, hints string) string {
	content := hints
	if errMsg != "" {
		content = theme.ErrorStyle.Render("✗ " + errMsg)
	} else if statusMsg != "" {
		content = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorAccent)).Render(statusMsg)
	}
	return theme.StatusBarStyle.Width(width).Render(content)
}

func HelpOverlay(width, height int) string {
	lines := []string{
		theme.TitleStyle.Render("KEYBOARD SHORTCUTS"),
		"",
		theme.SubtleStyle.Render("Navigation"),
		"  j/k / ↑↓   Move up/down",
		"  g/G         Top/bottom",
		"  1-6         Switch view",
		"  tab         Cycle panes",
		"",
		theme.SubtleStyle.Render("Search & Playback"),
		"  /           Search",
		"  enter       Tune station",
		"  space       Play/pause",
		"  x           Stop",
		"  m           Mute",
		"  -/=         Volume down/up",
		"",
		theme.SubtleStyle.Render("Station Actions"),
		"  f           Toggle favorite",
		"  r           Find related",
		"  o           Open detail",
		"  y           Copy stream URL",
		"  c           Add to crate",
		"  p           Save preset",
		"",
		theme.SubtleStyle.Render("History"),
		"  [/]         Prev/next in history",
		"",
		theme.SubtleStyle.Render("General"),
		"  ?           Toggle help",
		"  q / esc     Back/quit",
	}

	content := strings.Join(lines, "\n")
	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.ColorAccent)).
		Background(lipgloss.Color(theme.ColorSurface)).
		Foreground(lipgloss.Color(theme.ColorText)).
		Padding(1, 2).
		Width(width / 2).
		Height(height - 4)

	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		overlayStyle.Render(content))
}
