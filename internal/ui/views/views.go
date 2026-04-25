package views

import (
	"fmt"
	"strings"

	"github.com/cempack/radio-cli/internal/domain"
	"github.com/cempack/radio-cli/internal/theme"
	"github.com/cempack/radio-cli/internal/ui/components"
	"github.com/charmbracelet/lipgloss"
)

func Sidebar(width, height int, activeView int, playerState string) string {
	navItems := []struct {
		key  string
		name string
		id   int
	}{
		{"1", "Home", 0},
		{"2", "Discover", 1},
		{"3", "Favorites", 2},
		{"4", "History", 3},
		{"5", "Presets", 4},
		{"6", "Crates", 5},
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ColorAccent)).
		Bold(true).
		Padding(0, 1).
		Render("radiodrift")

	var navLines []string
	navLines = append(navLines, title)
	navLines = append(navLines, theme.SubtleStyle.Render(strings.Repeat("─", width-4)))
	navLines = append(navLines, "")

	for _, item := range navItems {
		label := fmt.Sprintf(" %s  %s", item.key, item.name)
		if item.id == activeView {
			navLines = append(navLines, lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.ColorAccent)).
				Bold(true).
				Width(width-4).
				Render(label))
		} else {
			navLines = append(navLines, lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.ColorText)).
				Width(width-4).
				Render(label))
		}
	}

	// fill remaining space
	usedLines := len(navLines)
	remaining := height - usedLines - 5
	for i := 0; i < remaining; i++ {
		navLines = append(navLines, "")
	}

	navLines = append(navLines, theme.SubtleStyle.Render(strings.Repeat("─", width-4)))
	navLines = append(navLines, theme.SubtleStyle.Render(" "+playerState))
	navLines = append(navLines, "")
	navLines = append(navLines, theme.SubtleStyle.Render(" ? for help"))

	content := strings.Join(navLines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func HomeView(width, height int, favorites []domain.Station, selectedIdx int) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("Welcome to radiodrift"))
	lines = append(lines, "")

	if len(favorites) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No favorites yet. Press / to search for stations."))
		lines = append(lines, "")
		lines = append(lines, theme.SubtleStyle.Render("Press 2 to discover random stations."))
		lines = append(lines, theme.SubtleStyle.Render("Press ? for help."))
	} else {
		lines = append(lines, theme.SubtleStyle.Render("Your Favorites:"))
		lines = append(lines, "")
		listWidth := width - 6
		for i, s := range favorites {
			row := components.StationRow(s, i == selectedIdx, listWidth)
			lines = append(lines, row)
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func DiscoverView(width, height int, stations []domain.Station, selectedIdx int, loading bool) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("Discover"))
	lines = append(lines, theme.SubtleStyle.Render("Press r to refresh random picks"))
	lines = append(lines, "")

	if loading {
		lines = append(lines, theme.SubtleStyle.Render("Loading stations..."))
	} else if len(stations) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No stations loaded. Press r to fetch random stations."))
	} else {
		listWidth := width - 6
		for i, s := range stations {
			row := components.StationRow(s, i == selectedIdx, listWidth)
			lines = append(lines, row)
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func SearchView(width, height int, stations []domain.Station, selectedIdx int, loading bool, inputView string) string {
	var lines []string
	lines = append(lines, inputView)
	lines = append(lines, "")

	if loading {
		lines = append(lines, theme.SubtleStyle.Render("Searching..."))
	} else if len(stations) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No results. Type to search."))
	} else {
		listWidth := width - 6
		for i, s := range stations {
			row := components.StationRow(s, i == selectedIdx, listWidth)
			lines = append(lines, row)
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func FavoritesView(width, height int, stations []domain.Station, selectedIdx int) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("Favorites"))
	lines = append(lines, "")

	if len(stations) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No favorites yet."))
		lines = append(lines, theme.SubtleStyle.Render("Press f on any station to add it."))
	} else {
		listWidth := width - 6
		for i, s := range stations {
			row := components.StationRow(s, i == selectedIdx, listWidth)
			lines = append(lines, row)
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func HistoryView(width, height int, history []domain.HistoryItem, stations map[string]domain.Station, selectedIdx int) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("History"))
	lines = append(lines, "")

	if len(history) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No listening history yet."))
	} else {
		for i, item := range history {
			s, ok := stations[item.StationID]
			name := item.StationID
			if ok {
				name = s.Name
			}
			timeStr := item.StartedAt.Format("Jan 02 15:04")
			label := fmt.Sprintf("%s  %s", timeStr, name)
			if i == selectedIdx {
				lines = append(lines, theme.SelectedRowStyle.Width(width-6).Render(label))
			} else {
				lines = append(lines, label)
			}
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func PresetsView(width, height int, presets []domain.Preset, selectedIdx int) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("Presets"))
	lines = append(lines, "")

	if len(presets) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No presets saved."))
		lines = append(lines, theme.SubtleStyle.Render("Press p on a station to save a preset."))
	} else {
		for i, p := range presets {
			label := fmt.Sprintf("[%s] %s", p.Key, p.Name)
			if i == selectedIdx {
				lines = append(lines, theme.SelectedRowStyle.Width(width-6).Render(label))
			} else {
				lines = append(lines, label)
			}
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func CratesView(width, height int, crates []domain.Crate, selectedIdx int) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("Crates"))
	lines = append(lines, "")

	if len(crates) == 0 {
		lines = append(lines, theme.SubtleStyle.Render("No crates yet."))
		lines = append(lines, theme.SubtleStyle.Render("Press c on a station to add it to a crate."))
	} else {
		for i, c := range crates {
			label := fmt.Sprintf("%s  (%d stations)", c.Name, len(c.StationIDs))
			if i == selectedIdx {
				lines = append(lines, theme.SelectedRowStyle.Width(width-6).Render(label))
			} else {
				lines = append(lines, label)
			}
		}
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}

func StationDetailPanel(s domain.Station, width, height int) string {
	var lines []string
	lines = append(lines, theme.TitleStyle.Render("Station Detail"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(s.Name))
	lines = append(lines, "")
	if s.Country != "" {
		lines = append(lines, fmt.Sprintf("Country:  %s", s.Country))
	}
	if s.State != "" {
		lines = append(lines, fmt.Sprintf("State:    %s", s.State))
	}
	if s.Language != "" {
		lines = append(lines, fmt.Sprintf("Language: %s", s.Language))
	}
	if s.Codec != "" {
		lines = append(lines, fmt.Sprintf("Codec:    %s", s.Codec))
	}
	if s.Bitrate > 0 {
		lines = append(lines, fmt.Sprintf("Bitrate:  %d kbps", s.Bitrate))
	}
	lines = append(lines, fmt.Sprintf("Votes:    %d", s.Votes))
	lines = append(lines, fmt.Sprintf("Clicks:   %d", s.ClickCount))
	if s.HomepageURL != "" {
		lines = append(lines, "")
		lines = append(lines, theme.SubtleStyle.Render("Homepage:"))
		lines = append(lines, theme.SubtleStyle.Render(s.HomepageURL))
	}
	if s.StreamURL != "" {
		lines = append(lines, "")
		lines = append(lines, theme.SubtleStyle.Render("Stream URL:"))
		lines = append(lines, theme.SubtleStyle.Render(s.StreamURL))
	}
	if len(s.Tags) > 0 {
		lines = append(lines, "")
		lines = append(lines, theme.SubtleStyle.Render("Tags:"))
		tagLine := ""
		for _, t := range s.Tags {
			tagLine += theme.BadgeStyle.Render(t) + " "
		}
		lines = append(lines, strings.TrimSpace(tagLine))
	}

	content := strings.Join(lines, "\n")
	return theme.PaneStyle.Width(width).Height(height).Render(content)
}
