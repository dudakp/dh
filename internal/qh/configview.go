package qh

import (
	"dh/pkg/executor"
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"reflect"
	"strings"
)

const (
	cpTemplatesFolder = iota
	cpDbConnectionString
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	executorService *SqlExecutorService
	focusIndex      int
	inputs          []textinput.Model
	cursorMode      cursor.Mode

	placeholderTagToValue map[string]string
}

func NewViewModel(executorService *SqlExecutorService) model {
	confType := reflect.TypeOf(executor.SqlExecutorConfig{})
	numFields := confType.NumField()
	m := model{
		executorService:       executorService,
		inputs:                make([]textinput.Model, numFields),
		placeholderTagToValue: map[string]string{},
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		if i == 0 {
			t.Focus()
		}

		field := confType.Field(i)
		placeholder := field.Tag.Get("placeholder")
		t.Placeholder = placeholder
		t.PromptStyle = focusedStyle
		t.TextStyle = focusedStyle

		m.inputs[i] = t
	}

	return m
}

func (r model) Init() tea.Cmd {
	return textinput.Blink
}

func (r model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return r, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && r.focusIndex == len(r.inputs) {
				r.updateConfig()
				return r, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				r.focusIndex--
			} else {
				r.focusIndex++
			}

			if r.focusIndex > len(r.inputs) {
				r.focusIndex = 0
			} else if r.focusIndex < 0 {
				r.focusIndex = len(r.inputs)
			}

			cmds := make([]tea.Cmd, len(r.inputs))
			for i := 0; i <= len(r.inputs)-1; i++ {
				if i == r.focusIndex {
					cmds[i] = r.inputs[i].Focus()
					r.inputs[i].PromptStyle = focusedStyle
					r.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				r.inputs[i].Blur()
				r.inputs[i].PromptStyle = noStyle
				r.inputs[i].TextStyle = noStyle
			}

			return r, tea.Batch(cmds...)
		}
	}
	cmd := r.updateInputs(msg)
	return r, cmd
}

func (r model) View() string {
	var b strings.Builder

	for i := range r.inputs {
		b.WriteString(r.inputs[i].View())
		if i < len(r.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if r.focusIndex == len(r.inputs) {
		button = &focusedButton
	}
	_, err := fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	if err != nil {
		return ""
	}

	return b.String()
}

func (r *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(r.inputs))
	for i := range r.inputs {
		r.inputs[i], cmds[i] = r.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (r model) updateConfig() {
	conf := executor.SqlExecutorConfig{}
	for i, input := range r.inputs {
		switch i {
		case cpTemplatesFolder:
			conf.TemplatesPath = input.Value()
			break
		case cpDbConnectionString:
			conf.DbConnectionString = input.Value()
			break
		}
		r.executorService.WriteConfig(conf)
	}
}
