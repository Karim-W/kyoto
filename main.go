package main

import (
	"github.com/karim-w/kyoto/cli"
	"github.com/karim-w/kyoto/spotify"
	"github.com/karim-w/kyoto/tui"
)

func main() {
	cli.Execute()
	client := spotify.SourceAuth()
	tui.Start(client)
}
