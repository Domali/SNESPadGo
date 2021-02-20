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
			watchdog <- "ping"
			ctrl.Update(sd)
			ctrl.DrawController(win)
		}
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

func serialWatchdog(watchdog <-chan string, data chan<- string) {
	for {
		select {
		case <-watchdog:
			for len(watchdog) > 0 {
				<-watchdog
			}
		default:
			data <- "ping"
		}
		time.Sleep(time.Millisecond * 33)
	}
}

func loadComConfig() string {
	dat, err := ioutil.ReadFile("./com.cfg")
	if err != nil {
		fmt.Println("Error reading COM port config.")
		os.Exit(1)
	}
	return string(dat)
}
