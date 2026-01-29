package ui

import (
	"github.com/MattiaPun/SubTUI/internal/api"
	"github.com/MattiaPun/SubTUI/internal/integration"
	"github.com/MattiaPun/SubTUI/internal/player"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	focusSearch = iota
	focusSidebar
	focusMain
	focusSong
	focusPlaylist = 90
)

const (
	viewList = iota
	viewQueue
	viewLogin = 99
)

const (
	filterSongs = iota
	filterAlbums
	filterArtist
)

const (
	displaySongs = iota
	displayAlbums
	displayArtist
)

const (
	LoopNone = 0
	LoopAll  = 1
	LoopOne  = 2
)

var albumTypes = []string{"Random", "Favorites", "Recently Added", "Recently Played", "Most Played"}

type Styles struct {
	Subtle    lipgloss.AdaptiveColor
	Highlight lipgloss.AdaptiveColor
	Special   lipgloss.AdaptiveColor
}

var Theme Styles

var (
	borderStyle       lipgloss.Style
	activeBorderStyle lipgloss.Style
	loginBoxStyle     lipgloss.Style
	loginHeaderStyle  lipgloss.Style
	loginHelpStyle    lipgloss.Style
	popupStyle        lipgloss.Style
)

// --- MODEL ---
type model struct {
	textInput    textinput.Model
	songs        []api.Song
	albums       []api.Album
	artists      []api.Artist
	playlists    []api.Playlist
	playerStatus player.PlayerStatus

	// Navigation State
	focus               int
	cursorMain          int
	cursorSide          int
	cursorAddToPlaylist int
	mainOffset          int

	// Window Dimensions
	width  int
	height int

	// View Mode
	viewMode        int
	filterMode      int
	displayMode     int
	displayModePrev int

	// App State
	err              error
	loading          bool
	lastPlayedSongID string
	scrobbled        bool
	loginErr         string
	discordRPC       bool
	notify           bool

	// Integrations
	dbusInstance    *integration.Instance
	discordInstance *integration.DiscordInstance

	// Queue System
	queue      []api.Song
	queueIndex int
	loopMode   int

	// Stars
	starredMap map[string]bool

	// Login State
	loginInputs []textinput.Model
	loginFocus  int

	// Input State
	lastKey string

	// Overlay States
	showHelp      bool
	showPlaylists bool
	helpModel     HelpModel
}

type HelpModel struct {
	Width  int
	Height int
}

type ContentModel struct {
	Content string
}

type BackgroundWrapper struct {
	RenderedView string
}

type loginResultMsg struct {
	err error
}

type songsResultMsg struct {
	songs []api.Song
}

type albumsResultMsg struct {
	albums []api.Album
}

type artistsResultMsg struct {
	artists []api.Artist
}

type playlistResultMsg struct {
	playlists []api.Playlist
}

type starredResultMsg struct {
	result *api.SearchResult3
}

type playQueueResultMsg struct {
	result *api.PlayQueue
}

type viewStarredSongsMsg *api.SearchResult3

type createShareMsg struct {
	url string
}

type errMsg struct {
	err error
}

type statusMsg player.PlayerStatus

type SetDBusMsg struct {
	Instance *integration.Instance
}

type SetDiscordMsg struct {
	Instance *integration.DiscordInstance
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search songs..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	startMode := viewList
	if api.AppConfig.Server.Username == "" || api.AppConfig.Server.Password == "" || api.AppConfig.Server.URL == "" {
		startMode = viewLogin
	}

	return model{
		textInput:           ti,
		songs:               []api.Song{},
		focus:               focusSearch,
		cursorMain:          0,
		cursorSide:          0,
		cursorAddToPlaylist: 0,
		viewMode:            startMode,
		filterMode:          filterSongs,
		displayMode:         displaySongs,
		starredMap:          make(map[string]bool),
		lastPlayedSongID:    "",
		loginInputs:         initialLoginInputs(),
		lastKey:             "",
		showHelp:            false,
		showPlaylists:       false,
		helpModel:           NewHelpModel(),
		discordRPC:          api.AppConfig.App.DiscordRPC,
		notify:              api.AppConfig.App.Notifications,
	}
}

func (m model) Init() tea.Cmd {
	if m.viewMode == viewList {
		return tea.Batch(
			textinput.Blink,
			attemptLoginCmd(),
		)
	}

	return textinput.Blink
}

func initialLoginInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "http(s)://music.example.com"
	inputs[0].Width = 30
	inputs[0].Focus()
	inputs[0].Prompt = "URL:      "

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "username"
	inputs[1].Width = 30
	inputs[1].Prompt = "Username: "

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "password"
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].Width = 30
	inputs[2].Prompt = "Password: "

	return inputs
}

func InitStyles() {
	Theme.Subtle = lipgloss.AdaptiveColor{Light: api.AppConfig.Theme.Subtle[0], Dark: api.AppConfig.Theme.Subtle[1]}
	Theme.Highlight = lipgloss.AdaptiveColor{Light: api.AppConfig.Theme.Highlight[0], Dark: api.AppConfig.Theme.Highlight[1]}
	Theme.Special = lipgloss.AdaptiveColor{Light: api.AppConfig.Theme.Special[0], Dark: api.AppConfig.Theme.Special[1]}

	// Global Borders
	borderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Subtle)

	// Focused Border (Brighter)
	activeBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Highlight)

	loginBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Highlight).
		Padding(1, 4).
		Align(lipgloss.Center)

	// The "Welcome" header
	loginHeaderStyle = lipgloss.NewStyle().
		Foreground(Theme.Special).
		Bold(true).
		MarginBottom(1)

	// The footer instruction
	loginHelpStyle = lipgloss.NewStyle().
		Foreground(Theme.Subtle).
		MarginTop(2)

	// The popup
	popupStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Highlight).
		Padding(1, 2)

}
