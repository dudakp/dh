package qh

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ResultModel struct {
	table table.Model
}

func NewResultModel(resultSet [][]string) ResultModel {
	if len(resultSet) < 2 {
		return ResultModel{}
	}
	res := ResultModel{}
	t := buildTable(resultSet)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	res.table = t
	return res
}

func (m ResultModel) Init() tea.Cmd { return nil }

func (m ResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
			// TODO: add relevant detail
			//tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ResultModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func buildTable(resultSet [][]string) table.Model {
	var columns []table.Column
	var rows []table.Row
	for i, row := range resultSet {
		if i == 0 {
			for _, value := range row {
				columns = append(columns, table.Column{Title: value, Width: 25})
			}
		} else {
			r := table.Row{}
			for j := range columns {
				r = append(r, row[j])
			}
			rows = append(rows, r)
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	return t
}
