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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	filepicker   filepicker.Model
	list         list.Model
	choice       string
	loader       spinner.Model
	selectedFile string
	quitting     bool
	err          error
	state        int
}

const (
	menuState = iota
	fpState
	loadingState
	successState
	errorState
)

type clearErrorMsg struct{}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
)

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
	l.Title = "Welcome to Peer to Peer File Sharing!"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	fp := filepicker.New()
	fp.Height = 10
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md", ".pdf"}
	fp.CurrentDirectory, _ = os.Getwd()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		list:       l,
		filepicker: fp,
		loader:     s,
		state:      menuState,
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
		case "tab":
			if m.state == menuState {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = string(i)
					m.state = fpState
				}
			}
			return m, nil
		case "enter":
			if m.state == fpState && m.selectedFile != "" {
				m.state = loadingState
				return m, m.loader.Tick
			}
		case "esc":

		}
	case clearErrorMsg:
		m.err = nil

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var lcmd, fcmd, locmd tea.Cmd
	m.list, lcmd = m.list.Update(msg)
	m.filepicker, fcmd = m.filepicker.Update(msg)
	m.loader, locmd = m.loader.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		m.state = successState
		return m, tea.Quit
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		m.state = errorState
		return m, tea.Batch(fcmd, clearErrorAfter(2*time.Second))
	}

	return m, tea.Batch(lcmd, fcmd, locmd)
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
			s.WriteString("Pick a file:")
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

	case loadingState:
		s.WriteString(fmt.Sprintf("%s Loading...", m.loader.View()))

	case successState:
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("green")).Render("Operation successful!"))

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
	case successState:
		return "Success"
	case errorState:
		return "Error"
	}
	return "Unknown"
}

func Start() {
	m := initialModel()
	tm, err := tea.NewProgram(&m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return
	}
	mm := tm.(model)
	fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
}
