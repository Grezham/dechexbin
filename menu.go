package main

import "fmt"

type MenuOption struct {
	text     string
	selected bool
	mode     int
}

func NewMenuOption(text string, mode int) MenuOption {
	return MenuOption{
		text:     text,
		selected: false,
		mode:     mode,
	}
}

/*
This Menu struct may become the basis for all menu pages within the app.
Possible functionality:

	Option lists
	Easy movement and selection of options
*/
type Menu struct {
	Name    string
	Options []MenuOption
	Cursor  string
	index   int
}

func (m *Menu) AddOption(text string) {
	newOption := MenuOption{text: text}
	m.Options = append(m.Options, newOption)
}

func (m *Menu) NextOption() {
	m.index++
	if m.index >= len(m.Options) {
		m.index = 0
	}
}

func (m *Menu) PrevOption() {
	m.index--
	if size := len(m.Options); m.index < 0 {
		m.index = size - 1
	}
}

// I Don't know why this doesn't work
// Selected can't be changed to false for some reason
func (m *Menu) Reset() {
	for _, opt := range m.Options {
		if opt.selected {
			opt.selected = false
		}
	}
	m.index = 0
}

func (m *Menu) Select() {
	m.Options[m.index].selected = true
}

func (m *Menu) SelectedChoiceMode() int {
	for _, c := range m.Options {
		if c.selected {
			return c.mode
		}
	}

	return -1
}

func (m Menu) Info() string {
	return fmt.Sprintf("%s\n%v\n%s\n%d\n", m.Name, m.Options, m.Cursor, m.index)

}
