package main

import (
	"fmt"
	"math/rand"

	"github.com/Craftman2868/go-libs/app"
	"github.com/Craftman2868/go-libs/event"
	"github.com/Craftman2868/go-libs/terminal"
)

const (
	MAX_DROP_COUNT     = 2000
	MAX_NEW_DROP_COUNT = 3
	DROP_CHARS         = "||.,'`"
	DROP_CHAR_COUNT    = int32(len(DROP_CHARS))
	NEW_DROP_MAX_Y     = 7
	MAX_FPS            = 24
)

var showFps = false

type RainApp struct {
	app.TerminalApp
	showFps      bool
	drops        []drop
	screenBuf    []byte
	termW, termH uint16
}

type drop struct {
	x, y  uint16
	char  byte
	speed uint16
}

func (ap *RainApp) InitRain() {
	ap.TerminalApp.InitTerminal()

	ap.On("key", ap.on_key)

	ap.On("update", func(event.Event) {
		ap.update()
	})
}

func (app *RainApp) on_key(ev_ event.Event) {
	ev := ev_.(terminal.KeyEvent)

	switch ev.Key {
	case 'C':
		if !ev.Ctrl() {
			break
		}
		fallthrough
	case terminal.KEY_ESC:
		fallthrough
	case 'Q':
		app.Stop()
	case 'F':
		app.showFps = !app.showFps
	}
}

func (app *RainApp) randDropChar() byte {
	n := rand.Int31n(DROP_CHAR_COUNT - 1)
	return DROP_CHARS[n]
}

func (app *RainApp) newDrop() {
	var d drop
	d.x = uint16(rand.Int31n(int32(app.termW - 1)))
	d.y = uint16(rand.Int31n(int32(min(app.termH, NEW_DROP_MAX_Y))))
	d.char = app.randDropChar()
	d.speed = uint16(rand.Int31n(3) + 1)
	app.drops = append(app.drops, d)
}

func (app *RainApp) initScreenBuf() {
	app.screenBuf = make([]byte, int(app.termW)*int(app.termH))
}

func (app *RainApp) clearScreenBuf() {
	for i := range app.screenBuf {
		app.screenBuf[i] = ' '
	}
}

func (app *RainApp) checkDropsPos() {
	for i := 0; i < len(app.drops); i++ {
		if app.drops[i].x >= app.termW || app.drops[i].y >= app.termH {
			app.drops[i] = app.drops[len(app.drops)-1]
			app.drops = app.drops[:len(app.drops)-1]
			i--
		}
	}
}

func (app *RainApp) writeScreenBuf(x, y int, data []byte) {
	p := y*int(app.termW) + x
	for i := range data {
		if p+i >= len(app.screenBuf) {
			break
		}
		app.screenBuf[p+i] = data[i]
	}
}

func (app *RainApp) update() {
	app.termW, app.termH = terminal.GetSize()
	app.termH--

	if len(app.screenBuf) != int(app.termW)*int(app.termH) {
		app.initScreenBuf()
		app.checkDropsPos()
	}

	app.clearScreenBuf()

	for i := 0; i < MAX_NEW_DROP_COUNT && len(app.drops) < MAX_DROP_COUNT; i++ {
		app.newDrop()
	}

	for i := 0; i < len(app.drops); i++ {
		d := &app.drops[i]
		app.screenBuf[int(d.y)*int(app.termW)+int(d.x)] = d.char
		d.y += d.speed

		if d.y >= app.termH {
			app.drops[i] = app.drops[len(app.drops)-1]
			app.drops = app.drops[:len(app.drops)-1]
			i--
		}
	}

	if app.showFps {
		app.writeScreenBuf(0, 0, fmt.Appendf(nil, "%.2f", app.GetTps()))
	}

	terminal.SetCursorHome()
	for i := range int(app.termH) {
		terminal.WritelnBytes(app.screenBuf[i*int(app.termW) : (i+1)*int(app.termW)])
	}
}

func main() {
	var app RainApp

	app.InitRain()

	app.Init(MAX_FPS)
	app.Run()
	app.Quit()
}
