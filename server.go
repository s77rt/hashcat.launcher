package hashcatlauncher

import (
	"embed"
	"log"
	"net"
	"net/http"
)

//go:embed frontend/hashcat.launcher/build
var fs embed.FS

func (a *App) NewServer() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	go http.Serve(ln, http.FileServer(http.FS(fs)))

	a.Server = ln
}
