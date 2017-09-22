package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// ActionFunc is a function that handles a button being clicked on a notification window.
type ActionFunc func( /* probably a gtk window or something? */ )

// DisplayNotificationMessage is a message sent to the GUI worker thread to display a notification
// with the given properties.
type DisplayNotificationMessage struct {
	// Title of the event.
	Title string
	// Reason specifies a description of why you're being notified.
	Reason string
	// Action1Label is the top button's label.
	Action1Label string
	// Action1 is the function that handles when the top button is clicked.
	Action1 ActionFunc
	// Action2Label is the bottom button's label.
	Action2Label string
	// Action2 is the function that handles when the bottom button is clicked.
	Action2 ActionFunc
}

func main() {
	ctx := context.Background()
	messages := make(chan *DisplayNotificationMessage)

	go guiWorker(ctx, messages)

	time.Sleep(1 * time.Second)

	messages <- &DisplayNotificationMessage{
		Title:  "Hello",
		Reason: "Test 1",
	}

	time.Sleep(1 * time.Second)

	messages <- &DisplayNotificationMessage{
		Title:  "World",
		Reason: "Test 2",
	}

	<-ctx.Done()
}

// guiWorker controls the state of the UI of the application in the background. We send messages to
// a separate thread to update the UI, allowing us to be a little more flexible with the creation of
// new windows over the course of the application's life.
func guiWorker(ctx context.Context, messages <-chan *DisplayNotificationMessage) {
	// Initialise GTK environment.
	gtk.Init(nil)

	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)

	go handleMessage(ctx, messages)

	// Start GTK.
	gtk.Main()
}

// handleMessage waits for messages to come in, and handles them. Most of the time this will
// facilitate creating new notification windows in another thread.
func handleMessage(ctx context.Context, messages <-chan *DisplayNotificationMessage) {
	for {
		select {
		case message := <-messages:
			glib.IdleAdd(displayNotification, message)
		case <-ctx.Done():
			gtk.MainQuit()
		}
	}
}

// displayNotification displays a notification window with the given parameters.
func displayNotification(message *DisplayNotificationMessage) {
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTypeHint(gdk.WINDOW_TYPE_HINT_NOTIFICATION)
	window.SetKeepAbove(true)
	window.SetResizable(false)
	window.SetSizeRequest(300, 70)
	window.SetSkipTaskbarHint(true)

	grid, _ := gtk.GridNew()
	grid.SetMarginBottom(10)
	grid.SetMarginTop(10)
	grid.SetMarginStart(10)
	grid.SetMarginEnd(10)
	grid.SetColumnSpacing(10)
	grid.SetRowSpacing(10)
	window.Add(grid)

	titleLabel, _ := gtk.LabelNew(fmt.Sprintf("<b>%s</b>", message.Title))
	titleLabel.SetUseMarkup(true)
	titleLabel.SetHAlign(gtk.ALIGN_START)
	titleLabel.SetHExpand(true)
	titleLabel.SetVExpand(true)

	reasonLabel, _ := gtk.LabelNew(message.Reason)
	reasonLabel.SetHAlign(gtk.ALIGN_START)
	reasonLabel.SetHExpand(true)
	reasonLabel.SetVExpand(true)

	// Top left. Bottom left. Top right. Bottom right.
	grid.Attach(titleLabel, 0, 0, 1, 1)
	grid.Attach(reasonLabel, 0, 1, 1, 1)

	window.Connect("destroy", func() {
		fmt.Println("Destroyed a window")
	})

	window.ShowAll()

	// We move after showing, because otherwise it'll show, and some WMs could float it after it's
	// been displayed. If that happens it will likely end up stranded in the center of the display.
	window.Move(20, 20)
}
