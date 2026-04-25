package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/cempack/radio-cli/internal/config"
	"github.com/cempack/radio-cli/internal/domain"
	"github.com/cempack/radio-cli/internal/keymap"
	"github.com/cempack/radio-cli/internal/player"
	"github.com/cempack/radio-cli/internal/radio"
	"github.com/cempack/radio-cli/internal/search"
	"github.com/cempack/radio-cli/internal/store"
	"github.com/cempack/radio-cli/internal/theme"
	"github.com/cempack/radio-cli/internal/ui/components"
	"github.com/cempack/radio-cli/internal/ui/views"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewID int

const (
	ViewHome ViewID = iota
	ViewDiscover
	ViewSearch
	ViewFavorites
	ViewHistory
	ViewPresets
	ViewCrates
	ViewSettings
)

type SearchResultMsg struct{ Stations []domain.Station }
type DiscoverResultMsg struct{ Stations []domain.Station }
type ErrorMsg struct{ Err error }
type StatusMsg struct{ Text string }
type PlayerStoppedMsg struct{}
type tickMsg time.Time

type AppModel struct {
	width, height int
	ready         bool
	activeView    ViewID

	searchInput   textinput.Model
	searchActive  bool
	searchResults []domain.Station
	selectedIdx   int

	currentStation *domain.Station
	player         *player.Player

	favorites []domain.Station
	history   []domain.HistoryItem
	presets   []domain.Preset
	crates    []domain.Crate

	discoverStations []domain.Station

	statusMsg string
	errMsg    string
	loading   bool
	showHelp  bool

	store       *store.Store
	radioClient *radio.Client
	cfg         domain.Config

	stationMap    map[string]domain.Station
	playbackStart time.Time
}

func New() *AppModel {
	ti := textinput.New()
	ti.Placeholder = "Search stations..."
	ti.CharLimit = 100

	st, err := store.New()
	if err != nil {
		st = nil
	}

	cfg := config.DefaultConfig()
	if st != nil {
		if loaded, err := st.LoadConfig(); err == nil && loaded.PlayerBackend != "" {
			cfg = loaded
		}
	}

	m := &AppModel{
		searchInput: ti,
		player:      player.New(),
		radioClient: radio.New(),
		store:       st,
		cfg:         cfg,
		stationMap:  make(map[string]domain.Station),
	}

	if st != nil {
		if favs, err := st.LoadFavorites(); err == nil {
			m.favorites = favs
			for _, s := range favs {
				m.stationMap[s.ID] = s
			}
		}
		if hist, err := st.LoadHistory(); err == nil {
			m.history = hist
		}
		if presets, err := st.LoadPresets(); err == nil {
			m.presets = presets
		}
		if crates, err := st.LoadCrates(); err == nil {
			m.crates = crates
		}
	}

	return m
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		m.fetchDiscover(),
	)
}

func (m *AppModel) fetchDiscover() tea.Cmd {
	return func() tea.Msg {
		stations, err := m.radioClient.Random(20)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return DiscoverResultMsg{Stations: stations}
	}
}

func (m *AppModel) searchAPI(query string) tea.Cmd {
	return func() tea.Msg {
		stations, err := m.radioClient.Search(radio.SearchParams{
			Name:  query,
			Limit: 50,
			Order: "votes",
		})
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SearchResultMsg{Stations: stations}
	}
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tickMsg:
		cmds = append(cmds, tick())

	case SearchResultMsg:
		m.searchResults = msg.Stations
		m.loading = false
		m.selectedIdx = 0

	case DiscoverResultMsg:
		m.discoverStations = msg.Stations
		m.loading = false
		if m.activeView == ViewDiscover {
			m.selectedIdx = 0
		}

	case ErrorMsg:
		m.errMsg = msg.Err.Error()
		m.loading = false

	case StatusMsg:
		m.statusMsg = msg.Text

	case tea.KeyMsg:
		m.errMsg = ""

		if m.showHelp {
			if key.Matches(msg, keymap.Default.Help) || key.Matches(msg, keymap.Default.Back) {
				m.showHelp = false
			}
			return m, tea.Batch(cmds...)
		}

		if m.searchActive {
			var inputCmd tea.Cmd
			m.searchInput, inputCmd = m.searchInput.Update(msg)
			cmds = append(cmds, inputCmd)

			switch {
			case key.Matches(msg, keymap.Default.Back):
				m.searchActive = false
				m.searchInput.Blur()
			case msg.String() == "enter":
				query := strings.TrimSpace(m.searchInput.Value())
				if query != "" {
					m.loading = true
					m.activeView = ViewSearch
					cmds = append(cmds, m.searchAPI(query))
				}
				m.searchActive = false
				m.searchInput.Blur()
			default:
				query := m.searchInput.Value()
				if query != "" && m.activeView == ViewSearch {
					// Filter across discover stations and favorites for broader real-time results
					pool := append(m.discoverStations, m.favorites...)
					m.searchResults = search.FilterStations(pool, search.SearchParams{Query: query})
				}
			}
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, keymap.Default.Back):
			if m.activeView == ViewSearch {
				m.activeView = ViewHome
				m.searchResults = nil
			} else {
				return m, tea.Quit
			}

		case key.Matches(msg, keymap.Default.Help):
			m.showHelp = !m.showHelp

		case key.Matches(msg, keymap.Default.Search):
			m.searchActive = true
			m.searchInput.Focus()
			m.activeView = ViewSearch

		case key.Matches(msg, keymap.Default.ViewHome):
			m.activeView = ViewHome
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.ViewDiscover):
			m.activeView = ViewDiscover
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.ViewFavorites):
			m.activeView = ViewFavorites
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.ViewHistory):
			m.activeView = ViewHistory
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.ViewPresets):
			m.activeView = ViewPresets
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.ViewCrates):
			m.activeView = ViewCrates
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.Up):
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case key.Matches(msg, keymap.Default.Down):
			max := m.currentListLen() - 1
			if m.selectedIdx < max {
				m.selectedIdx++
			}

		case key.Matches(msg, keymap.Default.Top):
			m.selectedIdx = 0

		case key.Matches(msg, keymap.Default.Bottom):
			m.selectedIdx = m.currentListLen() - 1

		case key.Matches(msg, keymap.Default.Tune):
			m.tuneSelected()

		case key.Matches(msg, keymap.Default.Stop):
			m.player.Stop()
			m.currentStation = nil
			m.statusMsg = "Stopped"

		case key.Matches(msg, keymap.Default.Favorite):
			m.toggleFavorite()
			m.saveFavorites()

		case key.Matches(msg, keymap.Default.Related):
			if m.activeView == ViewDiscover {
				m.loading = true
				cmds = append(cmds, m.fetchDiscover())
			}

		case key.Matches(msg, keymap.Default.PlayPause):
			if m.player.State() == player.StatePlaying {
				m.player.Stop()
				m.statusMsg = "Paused"
			} else if m.currentStation != nil {
				if err := m.player.Play(m.currentStation.StreamURL); err != nil {
					m.errMsg = err.Error()
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) currentList() []domain.Station {
	switch m.activeView {
	case ViewHome:
		return m.favorites
	case ViewDiscover:
		return m.discoverStations
	case ViewSearch:
		return m.searchResults
	case ViewFavorites:
		return m.favorites
	default:
		return nil
	}
}

func (m *AppModel) currentListLen() int {
	return len(m.currentList())
}

func (m *AppModel) tuneSelected() {
	list := m.currentList()
	if len(list) == 0 || m.selectedIdx >= len(list) {
		return
	}
	s := list[m.selectedIdx]
	m.currentStation = &s
	m.stationMap[s.ID] = s

	url := s.ResolvedURL
	if url == "" {
		url = s.StreamURL
	}

	if err := m.player.Play(url); err != nil {
		m.errMsg = err.Error()
		return
	}
	m.playbackStart = time.Now()
	m.statusMsg = fmt.Sprintf("Tuning: %s", s.Name)

	if m.cfg.SaveHistory {
		item := domain.HistoryItem{
			StationID: s.ID,
			StartedAt: m.playbackStart,
			Context:   fmt.Sprintf("%d", m.activeView),
		}
		m.history = append([]domain.HistoryItem{item}, m.history...)
		if len(m.history) > 100 {
			m.history = m.history[:100]
		}
		m.saveHistory()
	}
}

func (m *AppModel) toggleFavorite() {
	list := m.currentList()
	if len(list) == 0 || m.selectedIdx >= len(list) {
		return
	}
	s := list[m.selectedIdx]
	for i, f := range m.favorites {
		if f.ID == s.ID {
			m.favorites = append(m.favorites[:i], m.favorites[i+1:]...)
			m.statusMsg = fmt.Sprintf("Removed %s from favorites", s.Name)
			return
		}
	}
	m.favorites = append(m.favorites, s)
	m.statusMsg = fmt.Sprintf("Added %s to favorites", s.Name)
}

func (m *AppModel) saveFavorites() {
	if m.store != nil {
		_ = m.store.SaveFavorites(m.favorites)
	}
}

func (m *AppModel) saveHistory() {
	if m.store != nil {
		_ = m.store.SaveHistory(m.history)
	}
}

func (m *AppModel) View() string {
	if !m.ready {
		return "Loading radiodrift..."
	}

	if m.showHelp {
		return components.HelpOverlay(m.width, m.height)
	}

	statusBar := components.StatusBar(m.width, m.statusMsg, m.errMsg, m.hintLine())

	if m.width < 80 {
		return m.singleColumnView(statusBar)
	}
	if m.width < 120 {
		return m.twoColumnView(statusBar)
	}
	return m.threeColumnView(statusBar)
}

func (m *AppModel) hintLine() string {
	hints := []string{"/ search", "j/k nav", "enter tune", "f fav", "x stop", "? help", "q quit"}
	return theme.SubtleStyle.Render(strings.Join(hints, "  ·  "))
}

func (m *AppModel) mainContent(width, height int) string {
	switch m.activeView {
	case ViewHome:
		return views.HomeView(width, height, m.favorites, m.selectedIdx)
	case ViewDiscover:
		return views.DiscoverView(width, height, m.discoverStations, m.selectedIdx, m.loading)
	case ViewSearch:
		inputStr := m.searchInput.View()
		return views.SearchView(width, height, m.searchResults, m.selectedIdx, m.loading, inputStr)
	case ViewFavorites:
		return views.FavoritesView(width, height, m.favorites, m.selectedIdx)
	case ViewHistory:
		return views.HistoryView(width, height, m.history, m.stationMap, m.selectedIdx)
	case ViewPresets:
		return views.PresetsView(width, height, m.presets, m.selectedIdx)
	case ViewCrates:
		return views.CratesView(width, height, m.crates, m.selectedIdx)
	default:
		return views.HomeView(width, height, m.favorites, m.selectedIdx)
	}
}

func (m *AppModel) playerStateStr() string {
	switch m.player.State() {
	case player.StatePlaying:
		return "▶ playing"
	case player.StateConnecting:
		return "◌ connecting"
	case player.StateError:
		return "✗ error"
	default:
		return "■ stopped"
	}
}

func (m *AppModel) singleColumnView(statusBar string) string {
	contentH := m.height - lipgloss.Height(statusBar)
	main := m.mainContent(m.width, contentH)
	return lipgloss.JoinVertical(lipgloss.Left, main, statusBar)
}

func (m *AppModel) twoColumnView(statusBar string) string {
	rightW := m.width / 3
	leftW := m.width - rightW
	contentH := m.height - lipgloss.Height(statusBar)

	left := m.mainContent(leftW, contentH)
	right := components.NowPlayingPanel(m.currentStation, m.player, rightW, contentH)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	return lipgloss.JoinVertical(lipgloss.Left, body, statusBar)
}

func (m *AppModel) threeColumnView(statusBar string) string {
	sideW := m.width / 4
	rightW := m.width / 4
	centerW := m.width - sideW - rightW
	contentH := m.height - lipgloss.Height(statusBar)

	sidebar := views.Sidebar(sideW, contentH, int(m.activeView), m.playerStateStr())
	center := m.mainContent(centerW, contentH)
	right := components.NowPlayingPanel(m.currentStation, m.player, rightW, contentH)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, center, right)
	return lipgloss.JoinVertical(lipgloss.Left, body, statusBar)
}
