package main

import (
	"fmt"
	"log"
	"os"

	xcc "github.com/samanvp/ssxcc/pkg"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("%s input-file\n", args[0])
	}

	fd, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}

	inputProb, err := xcc.Builder(fd)
	inputProb.IsConsistent()

	inputProb.ProcessOption(2)
	inputProb.IsConsistent()
	fmt.Println(inputProb)
	inputProb.UndoOption(2)
	inputProb.IsConsistent()
	fmt.Println(inputProb)

	fmt.Println("main execution is done.")
}
