package xwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

// Selector
type Selector struct {
	widget.Select

	OnTapped func()
	OnChanged func(string)
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

func (s *Selector) SetSelected(string string) {
	s.Select.Selected = string
	s.Select.Refresh()
	if s.OnChanged != nil {
		s.OnChanged(string)
	}
}

// SelectorOption
type SelectorOption struct {
	widget.BaseWidget

	Text      string
	Alignment fyne.TextAlign
	TextStyle fyne.TextStyle

	Value string

	OnTapped func(string)
	OnTappedSecondary func(string)

	hovered  bool
}

func NewSelectorOptionWithStyle(text string, value string, alignment fyne.TextAlign, style fyne.TextStyle, tappedLeft func(string)) *SelectorOption {
	newselectoroption := &SelectorOption{}
	newselectoroption.ExtendBaseWidget(newselectoroption)
	newselectoroption.Text = text
	newselectoroption.Alignment = alignment
	newselectoroption.TextStyle = style
	newselectoroption.Value = value
	newselectoroption.OnTapped = tappedLeft
	newselectoroption.hovered = false
	return newselectoroption
}

func (so *SelectorOption) Tapped(*fyne.PointEvent) {
	if so.OnTapped != nil {
		so.OnTapped(so.Value)
	}
}

func (so *SelectorOption) TappedSecondary(*fyne.PointEvent) {
	if so.OnTappedSecondary != nil {
		so.OnTappedSecondary(so.Value)
	}
}

func (so *SelectorOption) MouseIn(*desktop.MouseEvent) {
	so.hovered = true
	so.Refresh()
}

func (so *SelectorOption) MouseOut() {
	so.hovered = false
	so.Refresh()
}

func (so *SelectorOption) MouseMoved(*desktop.MouseEvent) {
}

func (so *SelectorOption) CreateRenderer() fyne.WidgetRenderer {
	so.ExtendBaseWidget(so)
	label := widget.NewLabelWithStyle(so.Text, so.Alignment, so.TextStyle)
	background := canvas.NewRectangle(theme.BackgroundColor())
	return &SelectorOptionRenderer{
		so,
		background,
		label,
	}
}

type SelectorOptionRenderer struct {
	selectoroption *SelectorOption
	background *canvas.Rectangle
	label *widget.Label
}

func (h *SelectorOptionRenderer) Layout(size fyne.Size) {
	h.background.Resize(size)
}

func (h *SelectorOptionRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{h.background, h.label}
}

func (h *SelectorOptionRenderer) MinSize() (size fyne.Size) {
	size = h.label.MinSize()
	return
}

func (h *SelectorOptionRenderer) Refresh() {
	h.label.Refresh()
	if h.selectoroption.hovered {
		h.background.FillColor = theme.HoverColor()
	} else {
		h.background.FillColor = theme.BackgroundColor()
	}
	h.background.Refresh()
}

func (h *SelectorOptionRenderer) Destroy() {
}
