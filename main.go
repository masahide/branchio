package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/masahide/brachio/branchio"
)

func main() {
	file_a, err := os.Create("a.txt")
	if err != nil {
		log.Fatal(err)
	}
	file_b, err := os.Create("b.txt")
	if err != nil {
		log.Fatal(err)
	}
	a := branchio.NewBranchChannelWriter(file_a)
	b := branchio.NewBranchChannelWriter(file_b)

	branchWriter := branchio.NewWriter([]*branchio.BranchChannelWriter{a, b})

	io.Copy(branchWriter, os.Stdin)
	time.Sleep(2 * time.Second)

}
