package xwidget

import (
	"image/color"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/driver/desktop"
)

// Selector (Fake Selector)
type Selector struct {
	widget.Select

	OnTapped func()
}

func NewSelector(text string, tappedLeft func()) *Selector {
	newselector := &Selector{}
	newselector.ExtendBaseWidget(newselector)
	newselector.OnTapped = tappedLeft
	return newselector
}

func (s *Selector) Tapped(*fyne.PointEvent) {
	if s.OnTapped != nil {
		s.OnTapped()
	}
}

// SelectorOption (Fake SelectorOption)
type SelectorOption struct {
	*widget.Label

	Value string

	OnTapped func(string)

	hovered  bool
}

func NewSelectorOptionWithStyle(text string, value string, alignment fyne.TextAlign, style fyne.TextStyle, tappedLeft func(string)) *SelectorOption {
	return &SelectorOption{
		widget.NewLabelWithStyle(text, alignment, style),
		value,
		tappedLeft,
		false,
	}
}

func (so *SelectorOption) Tapped(*fyne.PointEvent) {
	if so.OnTapped != nil {
		so.OnTapped(so.Value)
	}
}

func (so *SelectorOption) TappedSecondary(*fyne.PointEvent) {
}

func (so *SelectorOption) MouseIn(*desktop.MouseEvent) {
	so.hovered = true
	canvas.Refresh(so)
}

func (so *SelectorOption) MouseOut() {
	so.hovered = false
	canvas.Refresh(so)
}

func (so *SelectorOption) MouseMoved(*desktop.MouseEvent) {
}

func (so *SelectorOption) CreateRenderer() fyne.WidgetRenderer {
	return &hoverLabelRenderer{so.Label.CreateRenderer(), so}
}

type hoverLabelRenderer struct {
	fyne.WidgetRenderer
	label *SelectorOption
}

func (h *hoverLabelRenderer) BackgroundColor() color.Color {
	if h.label.hovered {
		return theme.HoverColor()
	}

	return theme.BackgroundColor()
}
