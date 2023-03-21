package main

import "fmt"

type MenuOption struct {
	text string
	mode int
}

func NewMenuOption(text string, mode int) MenuOption {
	return MenuOption{
		text: text,
		mode: mode,
	}
}

type Menu struct {
	Name    string
	Options []MenuOption
	Cursor  string
	index   int
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

func (m *Menu) Select() int {
	return m.Options[m.index].mode
}

func (m Menu) Info() string {
	return fmt.Sprintf("%s\n%v\n%s\n%d\n", m.Name, m.Options, m.Cursor, m.index)

}
