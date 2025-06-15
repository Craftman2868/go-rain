package main

import (
	"fmt"
	"math/rand"

	"github.com/Craftman2868/go-libs/clock"
	"github.com/Craftman2868/go-libs/event"
	"github.com/Craftman2868/go-libs/terminal"
)

const (
	MAX_DROP_COUNT     = 2000
	MAX_NEW_DROP_COUNT = 4
	DROP_CHARS         = "||.,'`"
	DROP_CHAR_COUNT    = int32(len(DROP_CHARS))
	MAX_FPS            = 24
)

var showFps = false
var parser terminal.Parser
var handler event.BaseHandler
var appClock clock.Clock

func main() {
	terminal.Init()

	terminal.SetStyle("0")
	// terminal.SetStyle("0;1")
	appClock = clock.NewClock(MAX_FPS)

	handler.On("key", on_key)

	parser = terminal.NewParser(&handler)

	run()

	// terminal.SetStyle("0")

	terminal.Quit()
}

func on_key(ev_ event.Event) {
	ev := ev_.(terminal.KeyEvent)

	switch ev.Key {
	case 'C':
		if ev.Mod&terminal.MOD_CTRL == 0 {
			break
		}
		fallthrough
	case terminal.KEY_ESC:
		fallthrough
	case 'Q':
		running = false
	case 'F':
		showFps = !showFps
	}
}

type drop struct {
	x, y  uint16
	char  byte
	speed uint16
}

var running bool
var drops []drop

func run() {
	running = true

	for running {
		update()
		appClock.TickSleep()
	}
}

var termW, termH uint16

func randDropChar() byte {
	n := rand.Int31n(DROP_CHAR_COUNT - 1)

	return DROP_CHARS[n]
}

func newDrop() {
	var d drop

	d.x = uint16(rand.Int31n(int32(termW - 1)))
	d.y = uint16(rand.Int31n(10))
	d.char = randDropChar()
	d.speed = uint16(rand.Int31n(3) + 1)

	drops = append(drops, d)
}

var screenBuf []byte

func initScreenBuf() {
	screenBuf = make([]byte, int(termW)*int(termH))
}

func clearScreenBuf() {
	for i := range screenBuf {
		screenBuf[i] = ' '
	}
}

func checkDropsPos() {
	for i := 0; i < len(drops); i++ {
		if drops[i].x >= termW || drops[i].y >= termH {
			drops[i] = drops[len(drops)-1]
			drops = drops[:len(drops)-1]
			i--
		}
	}
}

func writeScreenBuf(x, y int, data []byte) {
	// screenBuf[y*int(termW)+x: math.Min(y*int(termW)+x+len(data), len(screenBuf))] = data[:math.Min(len(data), len(screenBuf)-(y*int(termW)+x))]
	p := y*int(termW) + x

	for i := range data {
		if p+i >= len(screenBuf) {
			break
		}
		screenBuf[p+i] = data[i]
	}
}

func update() {
	// Update initialization

	termW, termH = terminal.GetSize()
	termH--

	if termH < 10 {
		running = false
		return
	}

	if len(screenBuf) != int(termW)*int(termH) {
		initScreenBuf()
		checkDropsPos()
	}

	clearScreenBuf()

	// Handle inputs

	parser.HandleInputs()

	// Process rain

	for i := 0; i < MAX_NEW_DROP_COUNT && len(drops) < MAX_DROP_COUNT; i++ {
		newDrop()
	}

	for i := 0; i < len(drops); i++ {
		d := &drops[i]
		screenBuf[int(d.y)*int(termW)+int(d.x)] = d.char

		d.y += d.speed

		if d.y >= termH {
			drops[i] = drops[len(drops)-1]
			drops = drops[:len(drops)-1]
			i--
		}
	}

	// Show FPS

	if showFps {
		writeScreenBuf(0, 0, []byte(fmt.Sprint(appClock.GetTps())))
	}

	// Render screenBuf

	terminal.SetCursorHome()

	for i := range int(termH) {
		terminal.WritelnBytes(screenBuf[i*int(termW) : (i+1)*int(termW)])
	}
}

/*
main.update()
	/home/mateo/go/projects/rain/main.go:154 +0x3a9
main.run()
	/home/mateo/go/projects/rain/main.go:70 +0x1c
main.main()
	/home/mateo/go/projects/rain/main.go:34 +0x2a8
*/
