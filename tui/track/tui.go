package tui_track

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	common "github.com/hectorruiz-it/grabber/cmd"
)

type model struct {
	packages   []common.Repository
	index      int
	width      int
	height     int
	spinner    spinner.Model
	progress   progress.Model
	done       bool
	profile    string
	repository string
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func newModel(repositories []common.Repository, profile string) model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return model{
		packages: repositories,
		spinner:  s,
		progress: p,
		profile:  profile,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(downloadAndInstall(m.packages[m.index], m.profile), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case installedPkgMsg:
		pkg := m.packages[m.index].Path
		if m.index >= len(m.packages)-1 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Sequence(
				tea.Printf("%s %s", checkMark, pkg), // print the last success message
				tea.Quit,                            // exit the program
			)
		}

		// Update progress bar
		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.packages)))

		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, pkg), // print success message above our program
			downloadAndInstall(m.packages[m.index], m.profile), // download the next package
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	n := len(m.packages)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Added %d repositories to grabber.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := currentPkgNameStyle.Render(m.packages[m.index].Path)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Adding Repository " + pkgName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}

type installedPkgMsg string

func downloadAndInstall(pkg common.Repository, profile string) tea.Cmd {
	// This is where you'd do i/o stuff to download and install packages. In
	// our case we're just pausing for a moment to simulate the process.
	d := time.Millisecond //nolint:gosec
	return tea.Tick(d, func(t time.Time) tea.Msg {
		addRepositoryToConfig(profile, pkg)
		return installedPkgMsg(pkg.Path)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Track(repositories []common.Repository, profile string) {
	if _, err := tea.NewProgram(newModel(repositories, profile)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func addRepositoryToConfig(profile string, repository common.Repository) {
	config := common.ReadGrabberConfig()
	homeDir := common.GetHomeDirectory()
	// currentDir, err := os.Getwd()
	// common.CheckAndReturnError(err)

loop:
	for i := range config.Profiles {
		if config.Profiles[i].Profile == profile {
			for _, r := range config.Profiles[i].Repositories {
				if r == repository {
					continue loop
				}
			}
			config.Profiles[i].Repositories = append(config.Profiles[i].Repositories, repository)
		}
	}

	data, err := json.Marshal(config)
	common.CheckAndReturnError(err)
	err = os.WriteFile(homeDir+common.MAPPINGS_FILE_PATH, data, 0700)
	common.CheckAndReturnError(err)
}
