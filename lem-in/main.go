package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	filename := os.Args[1]

	graph, inputLines, err := ParseInput(filename)
	if err != nil {
		PrintError()
		return
	}

	bestPaths, err := FindOptimalPaths(graph)
	if err != nil {
		PrintError()
		return
	}
	fmt.Println(inputLines)
	Simulate(graph, bestPaths, inputLines)
}
