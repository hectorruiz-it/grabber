package progress_tui

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/zalando/go-keyring"
)

type model struct {
	repositories []Repository
	index        int
	width        int
	height       int
	spinner      spinner.Model
	done         bool
	doneRepos    map[string]bool // Track completed repositories by name
	ommitedRepos map[string]bool
	failedRepos  map[string]bool
}

type Repository struct {
	Path        string
	Profile     string
	ProfileType string
	Repository  string
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓").Bold(true)
	omitMark            = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).SetString("—").Bold(true)
	failedMark          = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("❌").Bold(true)

	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func newModel(migrationRepositories []Repository) model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	// Initialize done map
	doneRepos := make(map[string]bool)
	ommitedRepos := make(map[string]bool)
	failedRepos := make(map[string]bool)

	return model{
		repositories: migrationRepositories,
		spinner:      s,
		doneRepos:    doneRepos,
		ommitedRepos: ommitedRepos,
		failedRepos:  failedRepos,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, repo := range m.repositories {
		cmds = append(cmds, downloadAndInstall(repo))
	}
	cmds = append(cmds, m.spinner.Tick)

	return tea.Batch(cmds...)
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
		pkg := string(msg.repository)
		switch msg.status {
		case "cloned":
			m.doneRepos[pkg] = true
		case "ommited":
			m.ommitedRepos[pkg] = true
		case "failed":
			m.failedRepos[pkg] = true
		}

		m.index++

		if m.index >= len(m.repositories) {
			m.done = true
			return m, tea.Quit
		}

		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	for _, repo := range m.repositories {
		switch {
		case m.doneRepos[repo.Repository]:
			fmt.Fprintf(&b, "%s  %s\n", checkMark, repo.Repository)
		case m.ommitedRepos[repo.Repository]:
			fmt.Fprintf(&b, "%s  %s\n", omitMark, repo.Repository)
		case m.failedRepos[repo.Repository]:
			fmt.Fprintf(&b, "%s %s\n", failedMark, repo.Repository)
		default:
			spin := m.spinner.View() + " "
			pkgName := currentPkgNameStyle.Render(repo.Repository)
			fmt.Fprintf(&b, "%sCloning %s\n", spin, pkgName)
		}
	}

	if m.done {
		b.WriteString(doneStyle.Render(fmt.Sprintf("Migration complete! %d cloned, %d ommited and %d failed.\n", len(m.doneRepos), len(m.ommitedRepos), len(m.failedRepos))))
	}

	return b.String()
}

type installedPkgMsg struct {
	repository string
	status     string
}

func downloadAndInstall(repository Repository) tea.Cmd {
	return func() tea.Msg {
		// Perform the download and install asynchronously
		done := make(chan tea.Msg)

		go func() {
			sshRegex := regexp.MustCompile(`^git@`)
			httpsRegex := regexp.MustCompile(`^https://`)
			service := "grabber"

			InfoLog.Println("grabber: obtaining " + repository.Profile + "-profile entry from your keyring service.")

			password, err := keyring.Get(service, repository.Profile+"-profile")
			if err != nil {
				ErrorLog.Println("grabber:", err)
			}

			switch {
			case httpsRegex.MatchString(repository.Repository):

				_, err = git.PlainClone(repository.Path, false, &git.CloneOptions{
					Auth: &http.BasicAuth{
						Username: "git",
						Password: password,
					},
					URL: repository.Repository,
					// Progress: os.Stdout,
				})
				switch err {
				case nil:
					InfoLog.Println("grabber: succesfully cloned repository `" + repository.Repository + "`.")
					done <- installedPkgMsg{
						repository: repository.Repository,
						status:     "cloned",
					}
				case git.ErrRepositoryAlreadyExists:
					WarningLog.Println("grabber: repository `" + repository.Repository + "` already exists.")
					done <- installedPkgMsg{
						repository: repository.Repository,
						status:     "ommited",
					}
				default:
					ErrorLog.Println("grabber:", err)
					done <- installedPkgMsg{
						repository: repository.Repository,
						status:     "failed",
					}
				}

			case sshRegex.MatchString(repository.Repository):
				if repository.ProfileType == "ssh" {
					sshProfiles := common.ReadSshProfilesFile()
					privateKey, err := sshProfiles.Section(repository.Profile).GetKey("private_key")
					common.CheckAndReturnError(err)

					publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKey.Value(), password)
					common.CheckAndReturnError(err)

					_, err = git.PlainClone(repository.Path, false, &git.CloneOptions{
						Auth: publicKeys,
						// Progress: os.Stdout,
						URL: repository.Repository,
					})

					switch err {
					case nil:
						InfoLog.Println("grabber: succesfully cloned repository `" + repository.Repository + "`.")
						done <- installedPkgMsg{
							repository: repository.Repository,
							status:     "cloned",
						}
					case git.ErrRepositoryAlreadyExists:
						WarningLog.Println("grabber: repository `" + repository.Repository + "` already exists.")
						done <- installedPkgMsg{
							repository: repository.Repository,
							status:     "ommited",
						}
					default:
						ErrorLog.Println("grabber:", err)
						done <- installedPkgMsg{
							repository: repository.Repository,
							status:     "failed",
						}
					}
				} else {
					ErrorLog.Println("grabber: profile `" + repository.Profile + "` is not an ssh profile.")
				}
			}

			// Simulate some I/O work with a sleep
			// time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			// Send the completion message
			// done <- installedPkgMsg{
			// 	repository: repository.Repository,
			// 	status:     "failed",
			// }
		}()

		// Return the message received from the goroutine
		return <-done
	}
}

func ProgressTui(migrationRepositories []Repository) {
	dt := time.Now()
	logFilePath := common.GetHomeDirectory() + common.ROOT_DIR + "/migration-" + dt.Format("01-02-2006") + "_" + dt.Format("15:04:05")
	_, err := os.Create(logFilePath)
	common.CheckAndReturnError(err)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	common.CheckAndReturnError(err)
	defer logFile.Close()

	InfoLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	WarningLog = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLog = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)

	if _, err := tea.NewProgram(newModel(migrationRepositories)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
