package main

import (
	"log"

	"github.com/pocke/oshirase"
)

func main() {
	srv, err := oshirase.NewServer("melonnotifier", "onodera-punpun", "0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	n, err := newNotification(1920-56, 1200-(56*2), 56, "#EEEEEE",
		"#021B21", "/home/onodera/.fonts/cure.tff.bak", 11, 4)
	if err != nil {
		log.Fatal(err)
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
