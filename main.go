package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

/*type (
	errMsg error
)*/

const (
	binary      = 2
	decimal     = 10
	hexadecimal = 16
)

const (
	MainMenu = iota + 1
	SetMenu
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

type SetSettings struct {
	SetSize      int
	MaxRange     int
	QuestionType int
	AnswerType   int
}

type model struct {
	mode        int
	menu        Menu
	setMenu     Menu
	reviewMenu  Menu
	setSettings SetSettings
	CurrentMenu *Menu
	Set         QuestionSet
	index       int
	input       textinput.Model
}

func initialModel() model {
	//TODO: Move QuestionSet to Update loop after menu stuff is setup
	setSize, maxRange := 5, 100
	qSet := CreateQuestionSet(setSize, decimal, decimal, maxRange)
	input := textinput.New()
	input.Placeholder = "Answer Here..."
	input.Focus()

	nModel := model{
		mode:        MainMenu,
		menu:        CreateMainMenu(),
		setMenu:     CreateSetMenu(),
		reviewMenu:  CreateReviewMenu(),
		CurrentMenu: nil,
		setSettings: SetSettings{
			SetSize:      DefaultSetSize,
			MaxRange:     DefaultMaxRange,
			QuestionType: DefaultQuestionType,
			AnswerType:   DefaultAnswerType,
		},
		Set:   qSet,
		index: 0,
		input: input,
	}

	nModel.CurrentMenu = &nModel.menu
	return nModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit

		case "up":
			m.CurrentMenu.PrevOption()

		case "down":
			m.CurrentMenu.NextOption()

		case "enter":
			if m.mode == MainMenu {
				m.CurrentMenu.Select()
				switch m.CurrentMenu.SelectedChoiceMode() {
				case NewSet:
					m.mode = SetMenu
					m.CurrentMenu.Options[m.CurrentMenu.index].selected = false
					//m.CurrentMenu.Reset()
					m.CurrentMenu = &m.setMenu
				case Exit:
					return m, tea.Quit
				}
				return m, nil
			}

			if m.mode == SetMenu {
				m.mode = Quiz
				m.CurrentMenu.Options[m.CurrentMenu.index].selected = false
				//m.CurrentMenu.Reset()
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
					m.CurrentMenu = &m.reviewMenu
					return m, cmd
				} else {
					return m, cmd
				}
			}

			if m.mode == ReviewMenu {
				m.CurrentMenu.Select()
				m.Set.Reset()
				//m.CurrentMenu.Reset()
				switch m.CurrentMenu.SelectedChoiceMode() {
				case RestartSet:
					m.CurrentMenu.Options[m.CurrentMenu.index].selected = false
					//m.CurrentMenu.Reset()
					m.mode = Quiz
				case Exit:
					m.CurrentMenu.Options[m.CurrentMenu.index].selected = false
					//m.CurrentMenu.Reset()
					m.mode = MainMenu
					m.CurrentMenu = &m.menu
				}
			}
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {

	switch m.mode {
	case MainMenu:
		s := fmt.Sprintf("%s\n\n\n", m.CurrentMenu.Name)
		for i, opt := range m.CurrentMenu.Options {

			if i == m.CurrentMenu.index {
				s += fmt.Sprintf("[%s]", m.menu.Cursor)
			}
			s += fmt.Sprintf("\t%s\n", opt.text)
		}
		s += fmt.Sprintf("\n%s\n%d", m.CurrentMenu.Info(), m.mode)
		return s

	case SetMenu:
		s := fmt.Sprintf("%s\n\n\n", m.CurrentMenu.Name)
		for i, opt := range m.CurrentMenu.Options {
			if i == m.CurrentMenu.index {
				s += fmt.Sprintf("[%s]", m.menu.Cursor)
			}
			s += fmt.Sprintf("\t%s\n", opt.text)
		}
		s += fmt.Sprintf("\n%s\n%d", m.CurrentMenu.Info(), m.mode)
		return s

	case Quiz:
		s := fmt.Sprintf("Question %d\n\n", m.Set.GetQuestionNumber())
		for i := 0; i < m.Set.index; i++ {
			//TODO: REMOVE Results bool from formatting
			if m.Set.results[i] {
				s += fmt.Sprintf("%d %s :", i+1, DefaultCorrectMark)
			} else {
				s += fmt.Sprintf("%d %s :", i+1, DefaultWrongMark)
			}
		}
		s += fmt.Sprintf("\n%s %s\n\n%d", m.Set.GetCurrentQuestion(), m.input.View(), m.menu.SelectedChoiceMode())
		s += fmt.Sprintf("\n%s\n%d", m.CurrentMenu.Info(), m.mode)
		return s

	case ReviewMenu:
		s := fmt.Sprintf("%s\n\n\n", m.reviewMenu.Name)
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

		//Just need some space between results and options
		s += "\n"

		for i, opt := range m.CurrentMenu.Options {
			if i == m.CurrentMenu.index {
				s += fmt.Sprintf("[%s]", m.menu.Cursor)
			}
			s += fmt.Sprintf("\t%s\n", opt.text)
		}
		s += fmt.Sprintf("\n%s\n%d", m.CurrentMenu.Info(), m.mode)
		return s

	default:
		return "broke"
	}
}

// Menus
func CreateMainMenu() Menu {
	opt := make([]MenuOption, 2)
	opt[0] = NewMenuOption("New Set", NewSet)
	opt[1] = NewMenuOption("Exit", Exit)
	return Menu{
		Name:    "Main Menu",
		Options: opt,
		Cursor:  DefaultCursor,
		index:   0,
	}
}

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
	opt[1] = NewMenuOption("Exit to MainMenu", Exit)

	return Menu{
		Name:    "Review Menu",
		Options: opt,
		Cursor:  DefaultCursor,
		index:   0,
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
