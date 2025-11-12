package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alisdairrankine/frienvironment"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	app := frienvironment.Parse("GUI Test", `
// --- INTERRUPT SETUP ---
&ON-CLICK      // Get pointer to ON-CLICK function
2              // Interrupt type 2 (InterruptTypeMouseDown)
SYSCALL.INTERRUPT.REGISTER

// --- INITIAL DRAW ---
// Set a blue background
0 100 200 255 DRAW.CLEAR

// Draw a title
"Click to draw a rectangle!"
PRINT-STRING-TO-BUFFER
10 10 20 255 255 255 255 DRAW.TEXT

YIELD // Halt the VM and wait for interrupts

// --- INTERRUPT HANDLER ---
// (x y -- )
DEF ON-CLICK
    // When triggered, X and Y are on the stack
    10 10           // width, height
    255 255 0 255   // r, g, b, a (yellow)
    DRAW.RECT       // (x y w h r g b a -- )
    RET             // Return from interrupt

// Function to load a string into the buffer
DEF PRINT-STRING-TO-BUFFER
    out: BUF
    DUP 0 = IF
        DROP
        RET
    ELSE
        :out
    THEN
RET
`)

	vm := frienvironment.NewVM(app)

	if len(os.Args) == 2 && os.Args[1] == "docs" {
		fmt.Println(vm.Docs())
		return
	}

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	dq := &frienvironment.DrawQueue{}

	vm.AddWord("DRAW.CLEAR", "(r g b a -- ) Clears background", func() {

		a := uint8(vm.Stack.Pop())
		b := uint8(vm.Stack.Pop())
		g := uint8(vm.Stack.Pop())
		r := uint8(vm.Stack.Pop())

		dq.AddCommand(func() {
			rl.ClearBackground(rl.NewColor(r, g, b, a))
		})
	})

	vm.AddWord("DRAW.RECT", "(x y w h r g b a -- ) Draws a rectangle", func() {
		fmt.Println("rect", vm.Stack)
		a := uint8(vm.Stack.Pop())
		b := uint8(vm.Stack.Pop())
		g := uint8(vm.Stack.Pop())
		r := uint8(vm.Stack.Pop())
		h := int32(vm.Stack.Pop())
		w := int32(vm.Stack.Pop())
		y := int32(vm.Stack.Pop())
		x := int32(vm.Stack.Pop())

		dq.AddCommand(func() {
			rl.DrawRectangle(x, y, w, h, rl.NewColor(r, g, b, a))
		})
	})

	vm.AddWord("DRAW.TEXT", "(x y size r g b a -- ) Draws text from buffer", func() {

		a := uint8(vm.Stack.Pop())
		b := uint8(vm.Stack.Pop())
		g := uint8(vm.Stack.Pop())
		r := uint8(vm.Stack.Pop())
		size := int32(vm.Stack.Pop())
		y := int32(vm.Stack.Pop())
		x := int32(vm.Stack.Pop())

		text := string(vm.ReadFromBuffer())
		vm.ClearBuffer()

		dq.AddCommand(func() {
			rl.DrawText(text, x, y, size, rl.NewColor(r, g, b, a))
		})
	})

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
		// --- NEW: INPUT HANDLING ---
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			fmt.Print("CLICK")
			pos := rl.GetMousePosition()
			// Push Y, then X, so VM pops X then Y
			vm.Stack.Push(int(pos.X))
			vm.Stack.Push(int(pos.Y))
			vm.Interrupt(frienvironment.InterruptTypeMouseDown)
		}
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		dq.Run()

		rl.EndDrawing()
	}

}
