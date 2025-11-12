package main

import (
	"fmt"
	"io"
	"log"
	"os" // NEW
	"sync"

	"github.com/alisdairrankine/frienvironment"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// --- NEW: Windowing System ---

const (
	windowTitleHeight = 30 // Height of the title bar
)

var (
// Window management
)

// --- End Windowing System ---

var (
	drawQueue []func()
	queueLock sync.Mutex
)

func q(cmd func()) {
	queueLock.Lock()
	defer queueLock.Unlock()
	drawQueue = append(drawQueue, cmd)
}

func main() {

	f, err := os.Open("program.ali")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	// CHANGED: New program to demonstrate windowing
	app := frienvironment.Parse("Windowing Test", string(b))

	vm := frienvironment.NewVM(app)

	if len(os.Args) == 2 && os.Args[1] == "docs" {
		fmt.Println(vm.Docs())
		return
	}

	rl.InitWindow(1280, 720, "raylib [core] example - basic window")
	rl.ToggleFullscreen()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// --- VM WORD DEFINITIONS ---

	// This alert is no longer used, but we can leave it
	alert := ""
	vm.AddWord("SYSCALL.OUT", "DO NOT USE: debug function to output text", func() {
		alert = string(vm.ReadFromBuffer())
		vm.ClearBuffer()
	})

	wm := frienvironment.NewWindowManager()

	// CHANGED: All draw words are now window-relative

	vm.AddWord("WIN.DRAW.RECT", "(id x y w h r g b a -- ) Draws a rect in a window", func() {
		// Pop values from stack (in reverse order)
		a := uint8(vm.Stack.Pop())
		b := uint8(vm.Stack.Pop())
		g := uint8(vm.Stack.Pop())
		r := uint8(vm.Stack.Pop())
		h := int32(vm.Stack.Pop())
		w := int32(vm.Stack.Pop())
		y := int32(vm.Stack.Pop())
		x := int32(vm.Stack.Pop())
		id := vm.Stack.Pop()
		fmt.Println("draw rect")
		q(func() {
			win, ok := wm.WindowByID(id)
			if !ok {
				return // Window was closed
			}
			// Draw relative to window content area
			winX := win.X + x
			winY := win.Y + y + windowTitleHeight
			rl.DrawRectangle(winX, winY, w, h, rl.NewColor(r, g, b, a))
		})
	})

	vm.AddWord("WIN.DRAW.TEXT", "(id x y size r g b a -- ) Draws text from buffer in a window", func() {
		a := uint8(vm.Stack.Pop())
		b := uint8(vm.Stack.Pop())
		g := uint8(vm.Stack.Pop())
		r := uint8(vm.Stack.Pop())
		size := int32(vm.Stack.Pop())
		y := int32(vm.Stack.Pop())
		x := int32(vm.Stack.Pop())
		id := vm.Stack.Pop()

		text := string(vm.ReadFromBuffer())
		fmt.Println("text", text)
		vm.ClearBuffer()

		q(func() {
			win, ok := wm.WindowByID(id)
			if !ok {
				return
			}
			// Draw relative to window content area
			winX := win.X + x
			winY := win.Y + y + windowTitleHeight
			rl.DrawText(text, winX, winY, size, rl.NewColor(r, g, b, a))
		})
	})

	// NEW: Word to create a window
	vm.AddWord("WIN.CREATE", "(x y w h -- id) Creates a new window, title from buffer", func() {
		h := int32(vm.Stack.Pop())
		w := int32(vm.Stack.Pop())
		y := int32(vm.Stack.Pop())
		x := int32(vm.Stack.Pop())
		title := string(vm.ReadFromBuffer())
		vm.ClearBuffer()

		win := &frienvironment.Window{
			X:     x,
			Y:     y,
			W:     w,
			H:     h,
			Title: title,
		}
		wm.NewWindow(win)

		vm.Stack.Push(win.ID)
	})

	// NEW: Word to create a window
	vm.AddWord("WIN.DESTROY", "destroys a window", func() {
		id := vm.Stack.Pop()
		wm.DestroyWindow(id)

	})

	// We no longer use DRAW.CLEAR, but could implement WIN.CLEAR(id r g b a --)
	// which would draw a rect over the window's content area.
	// For now, we'll remove the old DRAW.CLEAR.

	// --------------------------

	go func() {
		err := vm.Run()
		if err != nil {
			log.Fatal(err)
		}
		if vm.State() != "waiting" {
			log.Fatalf("unexpected state: %s", vm.State())
		}
	}()

	for !rl.WindowShouldClose() {
		// --- CHANGED: Input & Drag Logic ---
		mousePos := rl.GetMousePosition()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			hitWindow := false
			// Get windows sorted top-to-bottom
			sortedWindows := wm.SortedWindows()
			// Loop from top to bottom
			for i := len(sortedWindows) - 1; i >= 0; i-- {
				win := sortedWindows[i]
				winRec := rl.NewRectangle(float32(win.X), float32(win.Y), float32(win.W), float32(win.H))

				if rl.CheckCollisionPointRec(mousePos, winRec) {
					// --- HIT A WINDOW ---
					hitWindow = true
					wm.BringToFront(win.ID)

					titleBarRec := rl.NewRectangle(float32(win.X), float32(win.Y), float32(win.W), windowTitleHeight)

					if rl.CheckCollisionPointRec(mousePos, titleBarRec) {
						// Start dragging
						wm.DraggingWindowID = win.ID
						wm.DragOffsetX = int32(mousePos.X) - win.X
						wm.DragOffsetY = int32(mousePos.Y) - win.Y
					} else {
						fmt.Println("clicked")
						// Clicked in content area, fire interrupt
						relX := int32(mousePos.X) - win.X
						relY := int32(mousePos.Y) - (win.Y + windowTitleHeight)
						vm.Stack.Push(int(relY))
						vm.Stack.Push(int(relX))
						vm.Stack.Push(win.ID)
						vm.Interrupt(frienvironment.InterruptTypeWindowMouseDown)
					}
					break // Stop after hitting the top-most window
				}
			}

			if !hitWindow {
				// --- HIT BACKGROUND ---
				vm.Stack.Push(int(mousePos.Y))
				vm.Stack.Push(int(mousePos.X))
				vm.Interrupt(frienvironment.InterruptTypeMouseDown)
			}
		}

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) && wm.DraggingWindowID != -1 {

			x := int32(mousePos.X) - wm.DragOffsetX
			y := int32(mousePos.Y) - wm.DragOffsetY
			wm.MoveWindow(wm.DraggingWindowID, x, y)
		}

		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			wm.DraggingWindowID = -1
		}
		// --- End Input Logic ---

		// --- CHANGED: Drawing Logic ---
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Welcome to AliOS", 5, 5, 30, rl.Black)

		// 1. Draw all window frames (sorted)
		sortedWindows := wm.SortedWindows()
		for _, win := range sortedWindows {
			// Draw window shadow
			rl.DrawRectangle(win.X+2, win.Y+2, win.W, win.H, rl.Gray)
			// Draw window background
			rl.DrawRectangle(win.X, win.Y, win.W, win.H, rl.LightGray)
			// Draw title bar
			rl.DrawRectangle(win.X, win.Y, win.W, windowTitleHeight, rl.Blue)
			rl.DrawText(win.Title, win.X+5, win.Y+3, 25, rl.White)
			// Draw outline
			rl.DrawRectangleLines(win.X, win.Y, win.W, win.H, rl.Black)
		}

		// 2. Process the draw queue (content)
		queueLock.Lock()
		for _, cmd := range drawQueue {
			cmd() // Execute the draw command
		}
		// drawQueue = []func(){} // Clear the queue
		queueLock.Unlock()

		// Draw alert if any
		if alert != "" {
			rl.DrawText(alert, 190, 200, 20, rl.Black)
		}

		rl.EndDrawing()
	}
}
