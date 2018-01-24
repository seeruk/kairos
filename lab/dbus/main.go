package main

import (
	"fmt"
	"log"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"github.com/skratchdot/open-golang/open"
	"strings"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	icon := "appointment-soon"

	hints := make(map[string]dbus.Variant)
	hints["urgency"] = dbus.MakeVariant([]byte{2})

	notification := notify.Notification{
		AppName:       "cnotifyd",
		ReplacesID:    0,
		AppIcon:       icon,
		Summary:       "Climbing at The Depot, Pudsey",
		Body:          "in 10 minutes",
		Actions:       []string{
			"open:https://www.google.co.uk/", "Open",
			"snooze", "Snooze",
		},
		Hints:         make(map[string]dbus.Variant),
	}

	notifier, err := notify.New(conn)
	if err != nil {
		panic(err)
	}
	defer notifier.Close()

	_, err = notifier.SendNotification(notification)
	if err != nil {
		panic(err)
	}

	actions := notifier.ActionInvoked()

	go func() {
		action := <-actions

		fmt.Printf("Action invoked: %v, with key %v\n", action.ID, action.ActionKey)

		if strings.HasPrefix(action.ActionKey, "open") {
			url := strings.SplitN(action.ActionKey, ":", 2)

			fmt.Println(url)
			fmt.Println(len(url))

			if len(url) != 2 {
				log.Println("wrong action format")
				return
			}

			log.Println(url[1])

			err := open.Run(url[1])
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	<-notifier.NotificationClosed()
}
