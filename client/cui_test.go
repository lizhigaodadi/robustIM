package client

import (
	"fmt"
	"github.com/rocket049/gocui"
	"testing"
)

func ShowGui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Printf("gocui init failed\n")
	}
	defer g.Close()

	v, err := g.SetView("viewname", 2, 2, 22, 7)
	if err != nil {
		if err != gocui.ErrUnknownView {
			// handle error
		}
		fmt.Fprintln(v, "This is a new view")
		// ...
	}
	fmt.Printf("%v", v)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {

	}

}

func TestShowGui(t *testing.T) {
	fmt.Printf("hello world\n")
	ShowGui()
}
