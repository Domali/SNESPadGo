package main

import (
	"bufio"
	"fmt"
	controller "github.com/domali/SNESPadGo/controller"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tarm/serial"
	"golang.org/x/image/colornames"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	ctrl := controller.NewController()
	cfg := pixelgl.WindowConfig{
		Title:  "SNESPad",
		Bounds: pixel.R(0, 0, ctrl.XSize, ctrl.YSize),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Couple variables used to calculate FPS of the program.
	frames := 0
	second := time.Tick(time.Second)

	watchdog := make(chan string, 10)
	data := make(chan string)
	go readSerialData(data)
	go serialWatchdog(watchdog, data)

	for !win.Closed() {
		win.Clear(colornames.Black)
		sd := <-data
		if sd != "ping" {
			// We received serial data so ping the watchdog
			watchdog <- "ping"
			// Update state of the buttons
			fmt.Println(sd)
			ctrl.Update(sd)
		}
		// Draw and update the controller every itteration even if no data was received.
		// If we only draw when there are serial updates the controller can flicker if the
		// console is doing something that prevents controller input for more than 33ms
		// (such as loading/saving astate in the MMX practice ROM).
		ctrl.DrawController(win)
		win.Update()
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		frames++
	}
	ctrl.SaveConfiguration()
}

func readSerialData(data chan<- string) {
	c := &serial.Config{Name: loadComConfig(), Baud: 57600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		data <- scanner.Text()
	}
}

// The run() function blocks on receiving data on the data channel. If the Teensy isn't
// sending data for whatever reason, such as the SNES being turned off, the program would
// block and freeze until the SNES is turned back on.  To prevent this we use a watchdog
// function that receives pings from the main program if it is updating. For every one
// check that the watchdog does it should receive around 2 pings. If the watchdog doesn't
// receive any pings it puts a ping on the data channel at 30 FPS. This will allow the main
// program to not block if there is no serial data being received.
func serialWatchdog(watchdog <-chan string, data chan<- string) {
	for {
		select {
		case <-watchdog:
			// Since the watchdog function checks at about half the speed we need to make sure we
			// remove any remaining pings on the watchdog channel. If we don't the main program will
			// fill up the buffer and block due to a race condition.
			for len(watchdog) > 0 {
				<-watchdog
			}
		default:
			data <- "ping"
		}
		time.Sleep(time.Millisecond * 33)
	}
}

// Instead of being com.cfg this could be some program configuration. Right now the
// only option is to select the COM port but we could easily add more configurations
// and change the name of this file.
func loadComConfig() string {
	dat, err := ioutil.ReadFile("./com.cfg")
	if err != nil {
		fmt.Println("Error reading COM port config.")
		os.Exit(1)
	}
	return string(dat)
}
