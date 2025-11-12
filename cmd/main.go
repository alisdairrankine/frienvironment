package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alisdairrankine/frienvironment"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	app := frienvironment.Parse("Hello World", `
"Loading... (1)"
PRINT-STRING
1000
SLEEP

"Loading... (2)"
PRINT-STRING
1000
SLEEP

"Loading... (3)"
PRINT-STRING
1000
SLEEP

"Hello World!"
PRINT-STRING
YIELD

DEF PRINT-STRING
    out: BUF
    DUP 0 = IF
        DROP
        SYSCALL.OUT
        RET
    ELSE
        :out
    THEN
RET`)

	vm := frienvironment.NewVM(app)

	if len(os.Args) == 2 && os.Args[1] == "docs" {
		fmt.Println(vm.Docs())
		return
	}

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	alert := ""

	vm.AddWord("SYSCALL.OUT", "DO NOT USE: debug function to output text", func() {
		alert = string(vm.ReadFromBuffer())
		vm.ClearBuffer()
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
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		if alert != "" {
			rl.DrawText(alert, 190, 200, 20, rl.Black)
		}
		rl.EndDrawing()
	}

}
