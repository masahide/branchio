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
	log.Printf("a:%v\n", branchWriter.CountWorkers())
	io.Copy(branchWriter, os.Stdin)
	for {
		if branchWriter.CountWorkers() == 0 {
			break
		}
		//time.Sleep(1 * time.Millisecond)
		time.Sleep(1 * time.Second)
		log.Printf("for:%v\n", branchWriter.CountWorkers())
	}

}
