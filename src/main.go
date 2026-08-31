package main

import (
	"log"
	"os"
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	doc         lipgloss.Style
	highlight   lipgloss.Style
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
	noTab       lipgloss.Style
	window      lipgloss.Style
	gapStyle    lipgloss.Style
}

type model struct {
	Width              int
	Height             int
	styles             *styles
	Tabs               []string
	TabContent         []string
	activeTab          int
	searchPackageInput textinput.Model
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func newStyles() *styles {
	darkMode := lipgloss.LightDark(true)

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	highlightColor := darkMode(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))

	s := new(styles)
	s.doc = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.activeTab = s.inactiveTab.Border(activeTabBorder, true)
	s.window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	s.gapStyle = lipgloss.NewStyle().Background(lipgloss.Color("235"))
	s.noTab = lipgloss.NewStyle().
		Foreground(highlightColor)

	return s
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = max(msg.Width-5, 5)
		m.Height = max(msg.Height-5, 5)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left", "right":
			if msg.String() == "left" {
				m.activeTab = max(m.activeTab-1, 0)
			} else {
				m.activeTab = (m.activeTab + 1) % len(m.Tabs)
			}
			if m.Tabs[m.activeTab] == "search" {
				// return m, m.searchPackageInput.Focus()
			} else {
				m.searchPackageInput.Blur()
			}
		}
	}

	if m.Tabs[m.activeTab] == "search" {
		m.searchPackageInput, cmd = m.searchPackageInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() tea.View {
	var cursor *tea.Cursor
	doc := strings.Builder{}
	s := m.styles

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = s.activeTab
		} else {
			style = s.inactiveTab
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┘"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)

	gapWidth := max(m.Width-lipgloss.Width(row), 0)
	emptyTabs := m.styles.noTab.Render(strings.Repeat("─", gapWidth))
	doc.WriteString(emptyTabs)

	doc.WriteString("\n")
	doc.WriteString(s.window.Width(m.Width).Render(m.TabContent[m.activeTab]))

	view := tea.NewView(s.doc.Render(doc.String()))
	view.AltScreen = true

	if !m.searchPackageInput.VirtualCursor() {
		cursor = m.searchPackageInput.Cursor()
		// cursor.Y += lipgloss.Height(doc.String())
	}
	view.Cursor = cursor

	return view
}

func (m model) viewPage() string {
	columns := []table.Column{
		{Title: "No.", Width: 10},
		{Title: "Package Name", Width: 40},
	}
	totalColumnWidth := 5
	for _, column := range columns {
		totalColumnWidth += column.Width
	}
	rows := []table.Row{
		{"1.", "Test"},
		{"2.", "Test"},
		{"3.", "Test"},
		{"4.", "Test"},
	}
	packageListTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
		table.WithWidth(totalColumnWidth),
	)

	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(false)
	style.Selected = lipgloss.NewStyle()
	packageListTable.SetStyles(style)

	return packageListTable.View()
}

func (m model) searchPage() string {
	textInput := textinput.New()
	textInput.Placeholder = "Package Name"
	textInput.CharLimit = 50
	textInput.SetWidth(50)
	textInput.SetVirtualCursor(false)

	m.searchPackageInput = textInput
	return "Package name: " + m.searchPackageInput.View()
}

func (m model) getTabContent() []string {
	var tabContent []string
	tabContent = append(tabContent, m.viewPage())
	tabContent = append(tabContent, m.searchPage())

	return tabContent
}

func createModel() model {
	tabs := []string{"view", "search"}

	m := model{Tabs: tabs, styles: newStyles()}
	m.TabContent = m.getTabContent()

	return m
}

func main() {
	m := createModel()
	program := tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatal("Error occured: %v", err)
		os.Exit(1)
	}
}
