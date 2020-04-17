package hashcatlauncher

import (
	"net/url"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func aboutScreen() fyne.CanvasObject {
	github_newissue_url, _ := url.Parse("https://github.com/s77rt/hashcat.launcher/issues/new/")
	hashcat_website_url, _ := url.Parse("https://hashcat.net/")
	hashc_url, _ := url.Parse("https://hashc.co.uk/")
	s77rt_email, _ := url.Parse("mailto:admin@abdelhafidh.com")
	return widget.NewVBox(
		widget.NewLabelWithStyle("What is hashcat.launcher:", fyne.TextAlignLeading, fyne.TextStyle{Bold:true}),
		widget.NewLabel("hashcat.launcher is a cross-platform app that run and control hashcat"),
		widget.NewLabel("it is designed to make it easier to use hashcat offering a friendly graphical user interface"),
		spacer(0, 10),
		widget.NewLabelWithStyle("Report a bug / Request a feature", fyne.TextAlignLeading, fyne.TextStyle{Bold:true}),
		widget.NewHyperlink("Continue to GitHub", github_newissue_url),
		spacer(0, 10),
		widget.NewLabelWithStyle("Useful links:", fyne.TextAlignLeading, fyne.TextStyle{Bold:true}),
		widget.NewHyperlink("hashcat Website", hashcat_website_url),
		widget.NewHyperlink("hashC Online Cracking Service", hashc_url),
		spacer(0, 10),
		widget.NewLabelWithStyle("License:", fyne.TextAlignLeading, fyne.TextStyle{Bold:true}),
		widget.NewLabel("MIT License"),
		widget.NewHBox(
			widget.NewLabel("Copyright (c) 2020 Abdelhafidh Belalia (s77rt)"),
			widget.NewHyperlink("<admin@abdelhafidh.com>", s77rt_email),
		),
	)
}
