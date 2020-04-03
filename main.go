package main

import (
	"log"

	"github.com/AndreKR/multiface"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/markbates/pkger"
	"github.com/onodera-punpun/oshirase"
)

var (
	runtime = pkger.Dir("/runtime")

	// Connection to the X server.
	X *xgbutil.XUtil

	// The multifont face that should be used.
	face *multiface.Face
)

func main() {
	// Initialize X.
	if err := initX(); err != nil {
		log.Fatalln(err)
	}

	// Initialize font face.
	if err := initFace(); err != nil {
		log.Fatalln(err)
	}

	// Initialize oshirase.
	srv, err := oshirase.NewServer("melonnotifier", "onodera-punpun", "0.0.1")
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize notification.
	n, err := initNotification(1920-56, 1200-(56*2), 56, xgraphics.BGRA{
		B: 238, G: 238, R: 238, A: 0xFF}, xgraphics.BGRA{B: 2, G: 27, R: 33,
		A: 0xFF}, 4)
	if err != nil {
		log.Fatalln(err)
	}

	ns := newNotifies()
	srv.OnNotify(func(o *oshirase.Notify) {
		ns.add(o)
		go n.show(o)
	})
	srv.OnCloseNotification(func(id uint32) bool {
		if err := ns.delete(id); err != nil {
			return false
		}
		return true
	})

	select {}
}
