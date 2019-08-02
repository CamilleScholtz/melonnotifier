package main

import (
	"github.com/AndreKR/multiface"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/zachomedia/go-bdf"
)

func initX() error {
	// Set up a connection to the X server.
	var err error
	X, err = xgbutil.NewConn()
	if err != nil {
		return err
	}

	// Run the main X event loop, this is used to catch events.
	go xevent.Main(X)

	return nil
}

func initFace() error {
	face = new(multiface.Face)

	fp, err := box.Find("fonts/cure.punpun.bdf")
	if err != nil {
		return err
	}
	f, err := bdf.Parse(fp)
	if err != nil {
		return err
	}
	face.AddFace(f.NewFace())

	fp, err = box.Find("fonts/kochi.small.bdf")
	if err != nil {
		return err
	}
	f, err = bdf.Parse(fp)
	if err != nil {
		return err
	}
	face.AddFace(f.NewFace())

	fp, err = box.Find("fonts/baekmuk.small.bdf")
	if err != nil {
		return err
	}
	f, err = bdf.Parse(fp)
	if err != nil {
		return err
	}
	face.AddFace(f.NewFace())

	return nil
}
