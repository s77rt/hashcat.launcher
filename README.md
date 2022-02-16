# hashcat.launcher
hashcat.launcher is a cross-platform app that run and control hashcat  
it is designed to make it easier to use hashcat offering a friendly graphical user interface

## Getting Started

### Requirements
 - [hashcat](https://hashcat.net/hashcat/)
 - Chromium based browser (Chrome, Edge, etc...)
 - `zenity`, `qarma` or `matedialog` (Linux only)
 - `osascript` (macOS only)

### Usage
 - Download a [release](https://github.com/s77rt/hashcat.launcher/releases)
 - Extract the archive
 - Run the executable

## Screenshots
![hashcat.launcher](/docs/screenshots/preview.gif?raw=true "hashcat.launcher")

## Building from source
requires [Go](https://go.dev/), [npm](https://www.npmjs.com/)
 - `git clone https://github.com/s77rt/hashcat.launcher.git`
 - `cd hashcat.launcher`
 - `make`
 - Executables can be found packaged in `bin` directory

## Changelog
Refer to [CHANGELOG.md](https://github.com/s77rt/hashcat.launcher/blob/master/docs/CHANGELOG.md)

___
[Report a bug / Request a feature](https://github.com/s77rt/hashcat.launcher/issues/new)