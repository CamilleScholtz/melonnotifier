package notify

import (
	"github.com/godbus/dbus"
)

// An Event with the `summary` of a notification.
// TODO: Can I somehow just use `chan string` instead of a type?
type Event struct {
	Notification chan string
}

func newEvent() *Event {
	return &Event{make(chan string)}
}

// Notify discards most "useless" info and just creates a notification event
// with the `summary`.
func (ev *Event) Notify(_ string, _ uint32, _ string, summary string, _ string,
	_ []string, _ map[string]dbus.Variant, _ int32) (uint32, *dbus.Error) {
	ev.Notification <- summary

	return 0, nil
}

// CloseNotification handles some freedesktop bullshit.
func (ev *Event) CloseNotification(_ uint32) *dbus.Error {
	return nil
}

// GetCapabilities handles some freedesktop bullshit.
func (ev *Event) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"actions", "body", "persistence"}, nil
}

// GetServerInformation handles some freedesktop bullshit.
func (ev *Event) GetServerInformation() (string, string, string, string,
	*dbus.Error) {
	return "melonnotifier", "onodera-punpun", "0.0.0", "1", nil
}
