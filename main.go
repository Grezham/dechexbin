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
	Quiz
)

const (
	NewSet = iota
	Exit
)

type model struct {
	mode  int
	menu  Menu
	Set   QuestionSet
	index []int
	input textinput.Model
}

func initialModel() model {
	//TODO: Will change to initiate mainmenu later
	setSize, maxRange := 5, 100
	qSet := CreateQuestionSet(setSize, decimal, decimal, maxRange)
	input := textinput.New()
	input.Placeholder = "Answer Here..."
	input.Focus()

	return model{
		mode:  MainMenu,
		menu:  CreateMainMenu(),
		Set:   qSet,
		index: []int{0, 0},
		input: input,
	}
}

//remove later

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
			m.menu.PrevOption()
		case "down":
			m.menu.NextOption()
		case "enter":
			//If at main menu start question
			//Will be changed to options
			/*
				New Set
				Settings
			*/
			if m.mode == MainMenu {
				m.mode = Quiz
				return m, nil
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
					m.mode = MainMenu
					//Resets the Set's done to false and index to 0
					m.Set.Reset()
					return m, cmd
				} else {
					return m, cmd
				}
			}
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.mode == MainMenu {
		s := "Welcome to Dechexbin!\n\n\n"
		for i, opt := range m.menu.Options {

			if i == m.menu.index {
				s += fmt.Sprintf("[%s]", m.menu.Cursor)
			}
			s += fmt.Sprintf("\t%s\n", opt.text)
		}
		return s
	}

	if m.mode == Quiz {
		s := fmt.Sprintf("Question %d\n\n", m.Set.GetQuestionNumber())
		for i := 0; i < m.Set.index; i++ {
			//TODO: REMOVE Results bool from formatting
			if m.Set.results[i] {
				s += fmt.Sprintf("● %v :", m.Set.results[i])
			} else {
				s += fmt.Sprintf("✖ %v :", m.Set.results[i])
			}
		}
		s += fmt.Sprintf("\n%s %s\n\n%s", m.Set.GetCurrentQuestion(), m.input.View(), "\nhelp stuff\n\n")
		return s
	}
	return "broke"
}

func CreateMainMenu() Menu {
	opt := make([]MenuOption, 2)
	opt[0] = MenuOption{text: "New Set", selected: true, mode: NewSet}
	opt[1] = MenuOption{text: "Exit", selected: false, mode: Exit}
	nMenu := Menu{
		Name:    "Main Menu",
		Options: opt,
		Cursor:  ">",
		index:   0,
	}

	return nMenu
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
