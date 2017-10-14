package notify

import (
	"fmt"

	"github.com/godbus/dbus"
)

// EventListener ...
func EventListener() (*Event, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	r, err := conn.RequestName("org.freedesktop.Notifications",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		return nil, err
	}
	if r != dbus.RequestNameReplyPrimaryOwner {
		return nil, fmt.Errorf("eavesdrop: Name already taken")
	}

	ev := newEvent()
	if err := conn.Export(ev, "/org/freedesktop/Notifications",
		"org.freedesktop.Notifications"); err != nil {
		return nil, err
	}

	return ev, nil
}
