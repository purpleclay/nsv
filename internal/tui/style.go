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

import "github.com/charmbracelet/lipgloss"

const (
	bottom       = "─"
	bottomLeft   = "└"
	bottomRight  = "┘"
	middle       = "┼"
	middleBottom = "┴"
	middleLeft   = "├"
	middleRight  = "┤"
	middleTop    = "┬"
	top          = "─"
	topLeft      = "┌"
	topRight     = "┐"
)

var (
	borderStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#2b0940"))
	hashStyle      = lipgloss.NewStyle().Background(lipgloss.Color("#1d1d1f")).Foreground(lipgloss.Color("#807d8a"))
	highlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("#bf31f7"))
	tagStyle       = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#3a1577"))
	feintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#807d8a"))
	chevron        = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#b769d6")).Render(">>")
	checkMark      = lipgloss.NewStyle().Foreground(lipgloss.Color("#139c20")).SetString("✓ ")
	diffMark       = lipgloss.NewStyle().Foreground(lipgloss.Color("#139c20")).SetString(" ↑↑")
	bullet         = feintStyle.Copy().SetString("> ")
	cell           = lipgloss.NewStyle()
)
