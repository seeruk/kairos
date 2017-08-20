package main

import (
	"fmt"
	"log"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	icon := "x-office-calendar"

	notification := notify.Notification{
		AppName:       "cnotifyd",
		ReplacesID:    0,
		AppIcon:       icon,
		Summary:       "test",
		Body:          "This is a test",
		Actions:       []string{"open", "Open"},
		Hints:         make(map[string]dbus.Variant),
		ExpireTimeout: 5000,
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

		fmt.Printf("Action invoked: %v, with key %v\n", action.Id, action.ActionKey)

		err := open.Run("https://www.elliotdwright.com")
		if err != nil {
			log.Println(err.Error())
		}
	}()

	<-notifier.NotificationClosed()
}
