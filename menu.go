package main

type MenuOption struct {
	text     string
	selected bool
	mode     int
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

func (m Menu) SelectedChoiceMode() int {
	for _, c := range m.Options {
		if c.selected {
			return c.mode
		}
	}

	return 0
}
