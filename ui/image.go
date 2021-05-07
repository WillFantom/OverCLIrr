package ui

import (
	"os"
	"runtime"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/browser"
	"github.com/willfantom/goverseerr"
	"golang.org/x/term"
)

func DisplayMediaPoster(posterPath string) {
	//Try in the CLI first
	err := openImageInTerminal(goverseerr.PosterPathBase + posterPath)
	if err != nil {
		Error("Can't open image in the command line, using browser")
		if bsrErr := openImageInBrowser(goverseerr.PosterPathBase + posterPath); bsrErr != nil {
			Fatal("Could not open image at all", err)
		}
	}
}

func openImageInBrowser(url string) error {
	return browser.OpenURL(url)
}

func openImageInTerminal(url string) error {
	var image *ansimage.ANSImage
	var x, y int
	var err error
	if term.IsTerminal(int(os.Stdout.Fd())) {
		x, y, err = term.GetSize(int(os.Stdout.Fd()))
	}
	if err != nil {
		x = 80
		y = 24
	}
	sm := ansimage.ScaleMode(2)
	dm := ansimage.DitheringMode(0)
	sfy, sfx := ansimage.BlockSizeY, ansimage.BlockSizeX
	if ansimage.DitheringMode(0) == ansimage.NoDithering {
		sfy, sfx = 2, 1
	}
	mc, _ := colorful.Hex("#000000")
	image, err = ansimage.NewScaledFromURL(url, sfy*y, sfx*x, mc, sm, dm)
	if err != nil {
		return err
	}
	if term.IsTerminal(int(os.Stdout.Fd())) {
		ansimage.ClearTerminal()
	}
	image.SetMaxProcs(runtime.NumCPU())
	image.DrawExt(false, false)
	return nil
}
