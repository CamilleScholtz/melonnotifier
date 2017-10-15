package main

import (
	"log"
	"time"

	"github.com/onodera-punpun/melonnotifier/notify"
)

func main() {
	ev, err := notify.EventListener()
	if err != nil {
		log.Fatal(err)
	}

	n, err := newNotification(1920-56, 1200-(56*2), 56, "#EEEEEE",
		"#021B21", "/home/onodera/.fonts/cure.tff.bak", 11)

	for {
		if err := n.draw(<-ev.Notification); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 4)
		n.undraw()
	}
}
