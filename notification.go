package main

import (
	"image"
	"os"
	"time"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/pocke/oshirase"
)

// Notification is a struct describing a notification popup.
type Notification struct {
	// Connection to the X server, the notify window, and the notify image.
	xu  *xgbutil.XUtil
	win *xwindow.Window
	img *xgraphics.Image

	// The X, Y and height of the notitication window.
	x, y, h int

	// The foreground and background colors of the notification window.
	bg, fg xgraphics.BGRA

	// The font and fontsize size that should be used.
	font *truetype.Font
	size float64

	// The time in seconds a notification is visible.
	time time.Duration

	// ID to check if we are killing the right window.
	ID uint32
}

func newNotification(x, y, h int, bg, fg, font string,
	size float64, time time.Duration) (n *Notification, err error) {
	n = new(Notification)

	// Set up a connection to the X server.
	n.xu, err = xgbutil.NewConn()
	if err != nil {
		return nil, err
	}

	// Run the main X event loop, this is uses to catch events.
	go xevent.Main(n.xu)

	// Create a window for the notification window. This window also listens to
	// button press events in order to respond to them.
	n.win, err = xwindow.Generate(n.xu)
	if err != nil {
		return nil, err
	}
	n.win.Create(n.xu.RootWin(), x, y, 600, h, xproto.CwBackPixel|
		xproto.CwEventMask, 0x000000, xproto.EventMaskButtonPress)

	// EWMH stuff to make the notification window visibile on all workspaces and
	// always be on top.
	if err := ewmh.WmWindowTypeSet(n.xu, n.win.Id, []string{
		"_NET_WM_WINDOW_TYPE_DOCK"}); err != nil {
		return nil, err
	}
	if err := ewmh.WmNameSet(n.xu, n.win.Id, "melonnotify"); err != nil {
		return nil, err
	}

	// Create the notification popup image.
	n.img = xgraphics.New(n.xu, image.Rect(0, 0, 600, h))
	n.img.XSurfaceSet(n.win.Id)

	// Set width and height of the notitication window.
	n.x = x
	n.y = y
	n.h = h

	// Convert foreground and background colors of the notification window from
	// HEX to BGRA.
	n.bg = hexToBGRA(bg)
	n.fg = hexToBGRA(fg)

	// Load font.
	// TODO: I don't *really* want to use `ttf` fonts but there doesn't seem to
	// be any `pcf` Go library at the moment. I have tried the plan9 fonts,
	// which do work, but honestly it's a pain in the ass (read: impossible) to
	// convert muh cure font.
	f, err := os.Open(font)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	n.font, err = xgraphics.ParseFont(f)
	if err != nil {
		return nil, err
	}
	n.size = size

	n.time = time

	// Listen to mouse events; close on click.
	xevent.ButtonPressFun(func(_ *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
		n.win.Unmap()
	}).Connect(n.xu, n.win.Id)

	return n, nil
}

func (n *Notification) show(o *oshirase.Notify) error {
	txt := "[" + o.Summary + "] " + o.Body

	// Calculate the required bar width coordinate.
	w, _ := xgraphics.Extents(n.font, n.size, txt)
	w += (2 * 24)
	if w > 600 {
		w = 600
	}

	// Color the background.
	n.img.For(func(cx, cy int) xgraphics.BGRA {
		return n.bg
	})

	// Draw the text.
	if _, _, err := n.img.Text(24, 20, n.fg, n.size, n.font, txt); err != nil {
		return err
	}

	// Make visible on all virtual desktops and map window.
	if err := ewmh.WmDesktopSet(n.xu, n.win.Id, 0xFFFFFFFF); err != nil {
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
		n.win.Unmap()
		n.ID = o.ID
	}

	return nil
}
