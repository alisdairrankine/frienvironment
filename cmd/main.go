package main

import (
	"io"
	"log"
	"os"

	"github.com/alisdairrankine/frienvironment"
)

func main() {

	f, err := os.Open("program.ali")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	app := frienvironment.Parse("Serial Test", string(b))

	vm := frienvironment.NewVM(app)

	if err := vm.Run(); err != nil {
		log.Fatal(err)
	}

}
