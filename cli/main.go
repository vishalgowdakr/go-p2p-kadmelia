package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	filepicker   filepicker.Model
	list         list.Model
	loader       spinner.Model
	viewport     viewport.Model
	selectedFile string
	choice       string
	quitting     bool
	err          error
	content      string
	ready        bool
	state        int
	substate     int
}

type workCompleteMsg struct{}

const (
	menuState = iota
	fpState
	loadingState
	errorState
	summaryState
)

const (
	uploadState = iota
	downloadState
)

type clearErrorMsg struct{}

const listHeight = 14

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

var (
	titleStyle_l      = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
)

func performWork() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return workCompleteMsg{}
	})
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func initialModel() model {
	items := []list.Item{
		item("Upload a file"),
		item("Download a file"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("Welcome to Peer to Peer File Sharing!")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle_l
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	fp := filepicker.New()
	fp.Height = 10
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md", ".pdf"}
	fp.CurrentDirectory, _ = os.Getwd()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	v := viewport.New(defaultWidth, listHeight)
	return model{
		list:       l,
		filepicker: fp,
		loader:     s,
		viewport:   v,
		state:      menuState,
		substate:   uploadState,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.filepicker.Init(),
		m.loader.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.state == menuState {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = string(i)
					if m.choice == "Download a file" {
						m.substate = downloadState
					}
					m.state = fpState
				}
				return m, nil
			} else if m.state == fpState && m.selectedFile != "" {
				m.state = loadingState
				return m, tea.Batch(
					m.loader.Tick,
					performWork(),
				)
			}
		case "esc":
			var cmd tea.Cmd
			m.filepicker, cmd = m.filepicker.Update(msg)
			return m, cmd

		}
	case clearErrorMsg:
		m.err = nil

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case workCompleteMsg:
		m.state = summaryState
		return m, nil
	}

	var lcmd, fcmd, locmd, viewportcmd tea.Cmd
	m.list, lcmd = m.list.Update(msg)
	m.filepicker, fcmd = m.filepicker.Update(msg)
	m.loader, locmd = m.loader.Update(msg)
	m.viewport, viewportcmd = m.viewport.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		m.state = loadingState
		return m, nil
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		m.state = errorState
		return m, tea.Batch(fcmd, clearErrorAfter(2*time.Second))
	}

	return m, tea.Batch(lcmd, fcmd, locmd, viewportcmd)
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString("Current state: " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(getState(m)) + "\n\n")
	switch m.state {
	case fpState:
		if m.quitting {
			return ""
		}
		s.WriteString("\n  ")
		if m.err != nil {
			s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
		} else if m.selectedFile == "" {
			if m.substate == uploadState {
				s.WriteString("Pick a file to upload:")
			} else {
				s.WriteString("Pick a '.torrent' file to download:")
			}
		} else {
			s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
		}
		s.WriteString("\n\n" + m.filepicker.View() + "\n")
		return s.String()

	case menuState:
		if m.choice != "" {
			s.WriteString(quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice)))
		} else if m.quitting {
			s.WriteString(quitTextStyle.Render("Goodbye!"))
		} else {
			s.WriteString(m.list.View())
		}
		s.WriteString("Your ID is: " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("1234567890") + "\n")

	case loadingState:
		if m.substate == uploadState {
			s.WriteString(fmt.Sprintf("%s Please wait while your file is being uploaded to our network...", m.loader.View()))
		} else {
			s.WriteString(fmt.Sprintf("%s Please wait while your file is being downloaded from our network", m.loader.View()))
		}

	case summaryState:
		if m.substate == uploadState {
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("green")).Render("File uploaded successfully!"))
		} else {
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("green")).Render("File downloaded successfully!"))
		}

	case errorState:
		s.WriteString(errorStyle.Render("An error occurred. Please try again."))
	}
	return s.String()
}

func getState(m model) string {
	switch m.state {
	case fpState:
		return "File Picker"
	case menuState:
		return "Menu"
	case loadingState:
		return "Loading"
	case summaryState:
		return "Summary"
	case errorState:
		return "Error"
	}
	return "Unknown"
}

func Start() {
	m := initialModel()
	_, err := tea.NewProgram(&m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return
	}
	/* mm := tm.(model)
	fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n") */
}
