package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
	"time"
)

const (
	WINDOW_WIDTH  = 1200
	WINDOW_HEIGHT = 1000

	SCROLL_LEN = 2048

	LOG_LINES_MAX = 20

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func setupWindow() *glfw.Window {
	if err := glfw.Init(); err != nil {
		closer.Fatalln(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Nuklear Demo", nil, nil)
	if err != nil {
		closer.Fatalln(err)
	}

	return win
}

func guiLoop(config *Config, state *State) {

	win := setupWindow()
	win.MakeContextCurrent()
	width, height := win.GetSize()

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 16, nil)
	idFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 32, nil)
	nk.NkFontStashEnd()
	if sansFont != nil {
		nk.NkStyleSetFont(ctx, sansFont.Handle())
	}

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	//dataMan := NewBmsDataManager()
	//dataMan.Set(BmsData{Packs: []Pack{pack}})

	fpsTicker := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			if win.ShouldClose() {
				close(exitC)
				continue
			}
			glfw.PollEvents()
			gfxMain(win, ctx, config, state, idFont)
		}
	}
}

func gfxMain(win *glfw.Window, ctx *nk.Context, config *Config, state *State, idFont *nk.Font) {

	width, height := win.GetSize()

	x := 0
	y := 0

	for i, _ := range state.PackData {
		if i == 6 {
			y += 300
			x = 0
		} else if i != 0 {
			x += 200
		}
		makePackFrame(ctx, config, state, idFont, i, float32(x), float32(y), 200, 300)
	}

	// Position in the bottom of the window
	makeLogFrame(ctx, state, 0, float32(height-400), float32(width), 400)
	//makeThresholdViewFrame(ctx, config, 0, 0, float32(width))

	// Render
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()

}

func b(v int32) bool {
	return v == 1
}

func flag(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
