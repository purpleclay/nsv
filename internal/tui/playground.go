/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	git "github.com/purpleclay/gitz"
	theme "github.com/purpleclay/lipgloss-theme"
)

type logEntryItem struct {
	description string
	title       string
}

func (i logEntryItem) Title() string       { return i.title }
func (i logEntryItem) Description() string { return i.description }
func (i logEntryItem) FilterValue() string { return i.title }

type gitLogLoadedMsg struct {
	log []list.Item
}

type PlaygroundOptions struct {
	GitC          *git.Client
	LatestTag     string
	Path          string
	VersionFormat string
}

type Playground struct {
	commit    textarea.Model
	commits   list.Model
	gitc      *git.Client
	latestTag string
	logPath   string
	pause     bool
}

func NewPlayground(opts PlaygroundOptions) *Playground {
	commits := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	commits.DisableQuitKeybindings()
	commits.SetShowHelp(false)
	commits.SetShowStatusBar(false)
	commits.SetShowTitle(false)

	commit := textarea.New()
	commit.SetHeight(8)

	return &Playground{
		commit:    commit,
		commits:   commits,
		gitc:      opts.GitC,
		latestTag: opts.LatestTag,
		logPath:   opts.Path,
	}
}

func (m *Playground) Init() tea.Cmd {
	return tea.Batch(m.gitLog)
}

func (m *Playground) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.commits.SetSize(msg.Width, msg.Height-13)
		m.commit.SetWidth(msg.Width)
	case gitLogLoadedMsg:
		cmds = append(cmds, m.commits.SetItems(msg.log))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Esc):
			m.commit.Blur()
			m.commit.Reset()
			m.pause = false
		case key.Matches(msg, Enter):
			// TODO: load commit message into the text area, and focus on the text area
			//
			logEntry := m.commits.SelectedItem().(logEntryItem)
			m.commit.SetCursor(0)
			m.commit.InsertString(logEntry.title)
			m.commit.SetCursor(0)
			cmds = append(cmds, m.commit.Focus())
			m.pause = true
		case key.Matches(msg, Quit):
			return m, tea.Quit
		}
	}

	if !m.pause {
		m.commits, cmd = m.commits.Update(msg)
		cmds = append(cmds, cmd)
	}

	// If not commit, update
	m.commit, cmd = m.commit.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Playground) View() string {
	// TODO: accurately calcualate the height based from this component
	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		theme.H2.Render("nsv"),
		theme.H4.Render("playground"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		"",
		m.commit.View(),
		"",
		m.commits.View(),
	)
}

func (m *Playground) gitLog() tea.Msg {
	var items []list.Item

	log, _ := m.gitc.Log(git.WithPaths(m.logPath), git.WithRefRange(git.HeadRef, m.latestTag))
	for _, l := range log.Commits {
		title := l.Message
		if idx := strings.Index(title, "\n"); idx > -1 {
			title = strings.TrimSpace(title[:idx])
		}

		items = append(items, logEntryItem{
			description: l.AbbrevHash,
			title:       title,
		})
	}

	return gitLogLoadedMsg{log: items}
}
