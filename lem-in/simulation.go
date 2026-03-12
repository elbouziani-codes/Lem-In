package main

import (
	"fmt"
	"strings"
)

// Ant tracks a single ant as it navigates its assigned path
type Ant struct {
	ID    int
	Path  []*Room
	Index int // -1 means at Start, len(Path) means End
}

// Simulate assigns ants to the optimal paths and tracks their turn-by-turn movements.
func Simulate(graph *Graph, paths [][]*Room, inputLines []string) {
	// 1. Output the input exactly as read, followed by a blank line
	for _, line := range inputLines {
		fmt.Println(line)
	}
	fmt.Println()

	// 2. Format paths: Add End room to the end to make simulating uniform
	fullPaths := make([][]*Room, len(paths))
	for i, p := range paths {
		fullPaths[i] = make([]*Room, len(p), len(p)+1)
		copy(fullPaths[i], p)
		fullPaths[i] = append(fullPaths[i], graph.End)
	}

	// 3. Pre-calculate optimal distribution of ants across the paths
	antsPerPath := make([]int, len(fullPaths))
	for i := 0; i < graph.NumAnts; i++ {
		bestP := -1
		bestTime := 1000000
		for p, path := range fullPaths {
			length := len(path) // Length from start to end (number of edges)
			time := length + antsPerPath[p]
			if time < bestTime {
				bestTime = time
				bestP = p
			}
		}
		antsPerPath[bestP]++
	}

	// 4. Simulate movements
	// activeAnts stores ants currently on a path.
	// We append new ants to the back, meaning activeAnts[0] is the oldest and furthest.
	activeAnts := make([][]*Ant, len(fullPaths))
	nextAntID := 1
	finishedAnts := 0

	for finishedAnts < graph.NumAnts {
		var moves []string

		for pIdx, p := range fullPaths {
			// First move ants already on the path
			for i := 0; i < len(activeAnts[pIdx]); i++ {
				ant := activeAnts[pIdx][i]
				ant.Index++
				room := p[ant.Index]

				moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, room.Name))

				if room == graph.End {
					finishedAnts++
				}
			}

			// Then spawn a new ant on this path space permitting
			if antsPerPath[pIdx] > 0 {
				ant := &Ant{
					ID:    nextAntID,
					Path:  p,
					Index: 0,
				}
				nextAntID++
				antsPerPath[pIdx]--

				room := p[0]
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, room.Name))

				if room == graph.End {
					finishedAnts++
				} else {
					// Pushed to back of slice. Will be moved next turn in the loops above!
					activeAnts[pIdx] = append(activeAnts[pIdx], ant)
				}
			}

			// Remove finished ants from active array to free memory and slice processing
			var remaining []*Ant
			for _, ant := range activeAnts[pIdx] {
				if ant.Path[ant.Index] != graph.End {
					remaining = append(remaining, ant)
				}
			}
			activeAnts[pIdx] = remaining
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}
