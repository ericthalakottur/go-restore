package main

import (
	"fmt"
	"os"
	"strings"

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
	Width      int
	Height     int
	styles     *styles
	Tabs       []string
	TabContent []string
	activeTab  int
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
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = max(msg.Width-5, 5)
		m.Height = max(msg.Height-5, 5)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case "right":
			m.activeTab = (m.activeTab + 1) % len(m.Tabs)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() tea.View {
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

	return view
}

func main() {
	tabs := []string{"view", "search"}
	tabContent := []string{"View Tab", "Search Tab"}
	m := model{Tabs: tabs, TabContent: tabContent, styles: newStyles()}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occured: %v", err)
		os.Exit(1)
	}
}
