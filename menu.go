package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
)

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

type InputValueInt struct {
	name      string
	value     int
	input     textinput.Model
	inputType int
}

func NewInputValueInt(name string, value int, inputType int) InputValueInt {

	return InputValueInt{
		name:      name,
		value:     value,
		input:     textinput.New(),
		inputType: inputType,
	}
}

func ValidSizeValue(value int) bool {
	return value > 0
}

func (iv *InputValueInt) ConvertInputValue() bool {
	if v, err := strconv.ParseInt(iv.input.Value(), 10, 64); err != nil {
		return false
	} else {
		if ValidSizeValue(int(v)) {
			iv.value = int(v)
			return true
		}
		return false
	}
}

func (iv *InputValueInt) Value() int {
	return iv.value
}

func (it *InputValueInt) SetValue(value int) {
	it.value = value
}

func (it *InputValueInt) Type() int {
	return it.inputType
}

func (it *InputValueInt) Name() string {
	return it.name
}

type InputToggle struct {
	name          string
	value         int
	toggleOptions []MenuOption
	toggleIndex   int
	inputType     int
}

func NewInputToggle(name string, value int, toggleOptions []MenuOption, inputType int) InputToggle {
	return InputToggle{
		name:          name,
		value:         value,
		toggleOptions: toggleOptions,
		toggleIndex:   0,
		inputType:     inputType,
	}
}

func (it *InputToggle) Value() int {
	return it.value
}

func (it *InputToggle) SetValue(value int) {
	it.value = value
}

func (it *InputToggle) Type() int {
	return it.inputType
}

func (it *InputToggle) Name() string {
	return it.name
}

func (it *InputToggle) ToggleNext() {
	it.toggleIndex++
	if it.toggleIndex >= len(it.toggleOptions) {
		it.toggleIndex = 0
	}
	it.value = it.toggleOptions[it.toggleIndex].mode
}

func (it *InputToggle) TogglePrev() {
	it.toggleIndex--
	if it.toggleIndex < 0 {
		it.toggleIndex = len(it.toggleOptions) - 1
	}
	it.value = it.toggleOptions[it.toggleIndex].mode
}

func (it *InputToggle) View() string {
	s := ""
	for i, t := range it.toggleOptions {
		if i == it.toggleIndex {
			s += fmt.Sprintf("{%s}", t.text)
		} else {
			s += t.text
		}
	}

	return s
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
