package hashcatlauncher

import (
	"embed"
	"net"
	"net/http"
)

//go:embed frontend/hashcat.launcher/build
var fs embed.FS

func (a *App) NewServer() error {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	go http.Serve(ln, http.FileServer(http.FS(fs)))

	a.Server = ln

	return nil
}
