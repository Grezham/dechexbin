package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	binary      = 2
	decimal     = 10
	hexadecimal = 16
)

const (
	SetMenu = iota + 1
	Quiz
	ReviewMenu
	NewSet
	RestartSet
	Exit
)

// Defaults
const (
	DefaultCursor           = "⇒"
	DefaultCorrectMark      = "●"
	DefaultWrongMark        = "✖"
	DefaultSetSize          = 10
	DefaultMaxRange         = 50
	DefaultQuestionType     = decimal
	DefaultAnswerType       = hexadecimal
	DefaultResultsLineLimit = 1
)

var (
	decimalOption     = NewMenuOption("Decimal", decimal)
	hexadecimalOption = NewMenuOption("Hexadecimal", hexadecimal)
	binaryOption      = NewMenuOption("Binary", binary)
)

var (
	SizeInput      = NewInputValueInt("Set Size", DefaultSetSize, 0)
	MaxRangeInput  = NewInputValueInt("Max", DefaultMaxRange, 0)
	questionToggle = NewInputToggle("Queston Type", decimal, []MenuOption{decimalOption, hexadecimalOption, binaryOption}, 0)
	answerToggle   = NewInputToggle("Answer Type", decimal, []MenuOption{decimalOption, hexadecimalOption, binaryOption}, 0)
)

type model struct {
	mode          int
	style         lipgloss.Style
	width, height int
	keys          keyMap
	help          help.Model
	setMenu       Menu
	reviewMenu    Menu
	CurrentMenu   Menu
	Set           *QuestionSet
	index         int
	input         textinput.Model
}

func initialModel() model {

	qSet := QuestionSet{} //CreateQuestionSet(setSize, decimal, decimal, maxRange)
	input := textinput.New()
	input.Placeholder = "Answer Here..."
	input.Focus()

	nModel := model{
		mode: SetMenu,
		style: lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center),
		keys:        keys,
		help:        help.New(),
		setMenu:     CreateSetMenu(),
		reviewMenu:  CreateReviewMenu(),
		CurrentMenu: Menu{},
		Set:         &qSet,
		index:       0,
		input:       input,
	}

	nModel.CurrentMenu = nModel.setMenu
	SizeInput.input.Focus()
	return nModel
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, cmd

		case key.Matches(msg, m.keys.Restart):
			if m.mode == Quiz {
				m.Set.Restart()
				m.input.Reset()
				return m, textinput.Blink
			} else if m.mode == ReviewMenu {
				m.Set.Restart()
				m.input.Reset()
				m.mode = Quiz
				return m, textinput.Blink
			}
			return m, cmd

		case key.Matches(msg, m.keys.Setting):
			if m.mode == Quiz {
				m.mode = SetMenu
				m.CurrentMenu = m.setMenu
			} else if m.mode == ReviewMenu {
				m.mode = SetMenu
				m.CurrentMenu = m.setMenu
			}
			return m, cmd

		case key.Matches(msg, m.keys.Up):
			m.CurrentMenu.PrevOption()
			if m.mode == SetMenu {
				switch m.CurrentMenu.index {
				case 0:
					SizeInput.input.Focus()
					MaxRangeInput.input.Blur()
				case 1:
					MaxRangeInput.input.Focus()
					SizeInput.input.Blur()
				default:
					MaxRangeInput.input.Blur()
					SizeInput.input.Blur()
				}
			}

		case key.Matches(msg, m.keys.Down):
			m.CurrentMenu.NextOption()
			if m.mode == SetMenu {
				switch m.CurrentMenu.index {
				case 0:
					SizeInput.input.Focus()
					MaxRangeInput.input.Blur()
				case 1:
					MaxRangeInput.input.Focus()
					SizeInput.input.Blur()
				default:
					MaxRangeInput.input.Blur()
					SizeInput.input.Blur()
				}

			}

		case key.Matches(msg, m.keys.Left):
			if m.mode == SetMenu {
				switch m.CurrentMenu.index {
				case 2:
					questionToggle.TogglePrev()
				case 3:
					answerToggle.TogglePrev()
				}
			}
		case key.Matches(msg, m.keys.Right):
			if m.mode == SetMenu {
				switch m.CurrentMenu.index {
				case 2:
					questionToggle.ToggleNext()
				case 3:
					answerToggle.ToggleNext()
				}
			}

		case key.Matches(msg, m.keys.Enter):
			if m.mode == SetMenu {
				//Remember to create something to randomize set if restarting from review
				validsize := SizeInput.ConvertInputValue()
				validRange := MaxRangeInput.ConvertInputValue()
				if validsize && validRange {
					m.Set = CreateQuestionSet(SizeInput.value, MaxRangeInput.value, questionToggle.Value(), answerToggle.Value())
					m.mode = Quiz
				}
				if !validsize {
					SizeInput.input.Reset()
					SizeInput.input.Placeholder = "Please enter valid whole decimal number"
				}
				if !validRange {
					MaxRangeInput.input.Reset()
					MaxRangeInput.input.Placeholder = "Please enter valid whole decimal number"
				}
				return m, textinput.Blink
			}

			if m.mode == Quiz {

				//Create GetAnswer() will check if correct
				m.Set.GetAnswer(m.input.Value())
				m.Set.CheckAnswer()

				//Clearing the inputbox
				m.input.Reset()

				//Move towards the next question in the set
				//The Set's done bool will be set is we've gone past the bounds of the Question array
				m.Set.NextQuestion()

				//Checks if Set's done bool is set to true
				if m.Set.isDone() {
					m.mode = ReviewMenu
					//Resets the Set's done to false and index to 0
					m.CurrentMenu = m.reviewMenu
					return m, cmd
				} else {
					return m, cmd
				}
			}

			if m.mode == ReviewMenu {
				m.Set.Reset()
				switch m.CurrentMenu.Select() {
				case RestartSet:
					m.mode = Quiz
					m.Set = CreateQuestionSet(SizeInput.value, MaxRangeInput.value, questionToggle.Value(), answerToggle.Value())
					return m, textinput.Blink
				case SetMenu:
					m.mode = SetMenu
					m.CurrentMenu = m.setMenu
				}

			}
		}
	}
	if m.mode == SetMenu {
		cmd = UpdateInputs(msg)
	} else if m.mode == Quiz {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {

	helpView := m.help.View(m.keys)
	s := ""
	switch m.mode {

	case SetMenu:
		s += fmt.Sprintf("%s\n\n\n", m.CurrentMenu.Name)
		s += fmt.Sprintf("%s\n", ViewSetupMenu(SizeInput.name, SizeInput.input.View(), 0, m.CurrentMenu.index, m.CurrentMenu.Cursor))
		s += fmt.Sprintf("%s\n", ViewSetupMenu(MaxRangeInput.name, MaxRangeInput.input.View(), 1, m.CurrentMenu.index, m.CurrentMenu.Cursor))
		s += fmt.Sprintf("%s\n", ViewSetupMenu(questionToggle.name, questionToggle.View(), 2, m.CurrentMenu.index, m.CurrentMenu.Cursor))
		s += fmt.Sprintf("%s\n", ViewSetupMenu(answerToggle.name, answerToggle.View(), 3, m.CurrentMenu.index, m.CurrentMenu.Cursor))

	case Quiz:
		s += fmt.Sprintf("Question %d\n\n", m.Set.GetQuestionNumber())
		for i := 0; i < m.Set.index; i++ {
			if m.Set.results[i] {
				s += fmt.Sprintf("%d %s :", i+1, DefaultCorrectMark)
			} else {
				s += fmt.Sprintf("%d %s :", i+1, DefaultWrongMark)
			}
		}
		s += fmt.Sprintf("\n%s %s\n\n", m.Set.GetCurrentQuestion(), m.input.View())

	case ReviewMenu:
		s += fmt.Sprintf("%s\n\n\n", m.reviewMenu.Name)
		values := 0
		for i, result := range m.Set.results {
			currentQuestion := m.Set.questions[i]
			//This currently is result per line thing is useless but I may need it later, so it stays for now
			if values >= DefaultResultsLineLimit {
				s += "\n"
				values = 0
			}

			if result {
				s += fmt.Sprintf("%s -- Q: %s | A: %s ", DefaultCorrectMark, currentQuestion.str, currentQuestion.answer)
			} else {
				s += fmt.Sprintf("%s -- Q: %s | A: %s ---- Want: %s", DefaultWrongMark, currentQuestion.str, currentQuestion.answer, currentQuestion.Want())
			}

			values++
		}

	default:
		s += "broke"
	}

	s += "\n\n" + helpView
	return m.style.Width(m.width).Height(m.height).Render(s)
}

func ViewSetupMenu(v1 string, v2 string, pos int, index int, cursor string) string {
	if pos == index {
		return fmt.Sprintf("%s %s %s", cursor, v1, v2)
	} else {
		return fmt.Sprintf("  %s %s", v1, v2)
	}
}

func UpdateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	SizeInput.input, cmds[0] = SizeInput.input.Update(msg)
	MaxRangeInput.input, cmds[1] = MaxRangeInput.input.Update(msg)

	return tea.Batch(cmds...)
}

// Menus

func CreateSetMenu() Menu {
	opt := make([]MenuOption, 4)
	opt[0] = NewMenuOption("Size", 0)
	opt[1] = NewMenuOption("Random Range", 0)
	opt[2] = NewMenuOption("Question Type", 0)
	opt[3] = NewMenuOption("Answer Type", 0)

	return Menu{
		Name:    "Setup Menu",
		Options: opt,
		Cursor:  DefaultCursor,
		index:   0,
	}
}

func CreateReviewMenu() Menu {
	opt := make([]MenuOption, 2)
	opt[0] = NewMenuOption("Restart Set", RestartSet)
	opt[1] = NewMenuOption("Set Settings", SetMenu)
	return Menu{
		Name:    "Review Menu",
		Options: opt,
		Cursor:  DefaultCursor,
		index:   0,
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
