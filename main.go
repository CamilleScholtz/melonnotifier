package main

import (
	"log"
	"time"

	"github.com/onodera-punpun/melonnotify/notify"
)

func main() {
	ev, err := notify.EventListener()
	if err != nil {
		log.Fatal(err)
	}

	n, err := newNotification(40, 40, 200, 56, "#EEEEEE", "#021B21",
		"/home/onodera/.fonts/cure.tff.bak", 11)

	for {
		n.draw(<-ev.Notification)
		time.Sleep(time.Second * 4)
		n.destroy()
	}
}
