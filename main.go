package main

import (
	"io"
	"log"
	"os"

	"github.com/masahide/branchio/lib"
)

func main() {
	file_in, err := os.Open("in.txt")
	if err != nil {
		log.Fatal(err)
	}
	file_a, err := os.Create("a.txt")
	if err != nil {
		log.Fatal(err)
	}
	file_b, err := os.Create("b.txt")
	if err != nil {
		log.Fatal(err)
	}
	/*
		a := branchio.NewBranchChannelWriter(file_a)
		b := branchio.NewBranchChannelWriter(file_b)

		branchWriter := branchio.NewWriter([]*branchio.BranchChannelWriter{a, b})
		log.Printf("a:%v\n", branchWriter.CountWorkers())
	*/

	branchWriter := lib.NewBranchWriter(100, file_a, file_b)
	io.Copy(branchWriter, file_in)
	branchWriter.Close()
	errs := branchWriter.Merge()
	log.Printf("errs:%# v\n", errs)
	file_a.Close()
	file_b.Close()

}
