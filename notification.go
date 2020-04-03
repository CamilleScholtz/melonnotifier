package main

import (
	"fmt"
	"image"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/onodera-punpun/oshirase"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Notification is a struct describing a notification popup.
type Notification struct {
	// The notify window, and the notify image.
	win *xwindow.Window
	img *xgraphics.Image

	// The X, Y and height of the notitication window.
	x, y, h int

	// The background color.
	bg xgraphics.BGRA

	// Drawer
	drawer *font.Drawer

	// The time in seconds a notification is visible.
	time time.Duration

	// ID to check if we are killing the right window.
	ID uint32
}

func initNotification(x, y, h int, bg, fg xgraphics.BGRA, time time.Duration) (
	n *Notification, err error) {
	n = new(Notification)

	// Run the main X event loop, this is uses to catch events.
	go xevent.Main(X)

	// Create a window for the notification window. This window also listens to
	// button press events in order to respond to them.
	n.win, err = xwindow.Generate(X)
	if err != nil {
		return nil, err
	}
	n.win.Create(X.RootWin(), x, y, 600, h, xproto.CwBackPixel|xproto.
		CwEventMask, 0x000000, xproto.EventMaskButtonPress)

	// EWMH stuff to make the window behave like an actual notification.
	if err := initEWMH(n.win.Id); err != nil {
		return nil, err
	}

	// Create the notification popup image.
	n.img = xgraphics.New(X, image.Rect(0, 0, 600, h))
	if err := n.img.XSurfaceSet(n.win.Id); err != nil {
		return nil, err
	}

	// Set width and height of the notitication window.
	n.x = x
	n.y = y
	n.h = h

	// Set background color.
	n.bg = bg

	// Set bar font face.
	n.drawer = &font.Drawer{
		Dst:  n.img,
		Src:  image.NewUniform(fg),
		Face: face,
	}

	// Set the display time.
	n.time = time

	return n, nil
}

func (n *Notification) show(o *oshirase.Notify) error {
	txt := "[" + o.Summary + "] " + o.Body

	// Calculate the required bar width coordinate.
	w := n.drawer.MeasureString(txt).Ceil()
	w += (2 * 24)
	if w > 600 {
		w = 600
	}

	// Color the background.
	n.img.For(func(cx, cy int) xgraphics.BGRA {
		return n.bg
	})

	// Draw the text.
	n.drawer.Dot = fixed.P(24, 32)
	n.drawer.DrawString(txt)

	// Listen to mouse events; close on click.
	xevent.ButtonPressFun(func(_ *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
		fmt.Println(o)
	}).Connect(X, n.win.Id)

	// Make visible on all virtual desktops and map window.
	if err := ewmh.WmDesktopSet(X, n.win.Id, ^uint(0)); err != nil {
		return err
	}
	n.win.Map()
	n.win.MoveResize(n.x-w, n.y, w, n.h)

	// Draw and paint image on window.
	n.img.XDraw()
	n.img.XPaint(n.win.Id)

	// Undraw notification.
	n.ID = o.ID
	time.Sleep(time.Second * n.time)
	if n.ID == o.ID {
		n.win.Detach()
		n.win.Unmap()
		n.ID = o.ID
	}

	return nil
}
