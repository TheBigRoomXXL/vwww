package main

import (
	"flag"
	"log"
)

func main() {
	numberOfPages := flag.Int("pages", 10_000, "number of pages")
	size := flag.Int("size", 0, "response size in KB")
	delay := flag.Int("delay", 100, "delay in ms")
	seed := flag.Int("seed", 42, "an int")

	flag.Parse()

	vwww := NewVWWW(*numberOfPages, *size, int64(*seed), *delay)
	log.Fatal(vwww.Serve(8080))
}
