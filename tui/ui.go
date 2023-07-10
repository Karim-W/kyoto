package tui

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karim-w/kyoto/spotify"
)

func Start(
	spot *spotify.Client,
) {
	m := model{
		spotifyClient: spot,
	}
	curr, err := m.spotifyClient.GetCurrentlyPlayingSong()
	if err != nil {
		panic(err)
	}
	if curr.Item.Name != m.song || curr.Item.Album.Name != m.album {
		m.album = curr.Item.Album.Name
		m.song = curr.Item.Name
		m.artist = curr.Item.Artists[0].Name
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// A model can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.
type model struct {
	spotifyClient *spotify.Client
	song          string
	artist        string
	album         string
}

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m model) Init() tea.Cmd {
	return tick
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		curr, err := m.spotifyClient.GetCurrentlyPlayingSong()
		if err != nil {
			return m, tick
		}
		if curr.Item.Name != m.song || curr.Item.Album.Name != m.album {
			m.album = curr.Item.Album.Name
			m.song = curr.Item.Name
			m.artist = curr.Item.Artists[0].Name
		}
		return m, tick
	}
	return m, nil
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m model) View() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		PaddingRight(4).
		PaddingBottom(2).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	str := fmt.Sprintf("%s - %s\n %s", m.song, m.album, m.artist)
	return style.Render(str)
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(3 * time.Second)
	return tickMsg{}
}
