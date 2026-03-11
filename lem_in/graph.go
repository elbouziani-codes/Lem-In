package lem_in

import (
	"log"
	"sort"
)

func GraphRoomsAndLinkes() {
	G.Network = make(map[string][]*Rooms)
	for i := range G.Links {
		from := G.Links[i].From
		to := G.Links[i].To

		G.Network[from.Name] = append(G.Network[from.Name], to)
		G.Network[to.Name] = append(G.Network[to.Name], from)
	}
}

// state is one BFS queue entry: the current room + the path taken so far + visited rooms.
type state struct {
	room    *Rooms
	path    []*Rooms
	visited map[string]bool
}

// bfsFindAllPaths finds every simple path from start to end using BFS.
func bfsFindAllPaths() [][]*Rooms {
	var allPaths [][]*Rooms

	initVisited := map[string]bool{G.RmStar.Name: true}
	queue := []state{{room: G.RmStar, path: []*Rooms{G.RmStar}, visited: initVisited}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, next := range G.Network[cur.room.Name] {
			if cur.visited[next.Name] {
				continue
			}

			newPath := make([]*Rooms, len(cur.path), len(cur.path)+1)
			copy(newPath, cur.path)
			newPath = append(newPath, next)

			if next == G.RmEnd {
				allPaths = append(allPaths, newPath)
				continue
			}

			newVisited := make(map[string]bool, len(cur.visited)+1)
			for k, v := range cur.visited {
				newVisited[k] = v
			}
			newVisited[next.Name] = true

			queue = append(queue, state{room: next, path: newPath, visited: newVisited})
		}
	}
	return allPaths
}

// pathsOverlap returns true if two paths share any intermediate room (excluding start and end).
func pathsOverlap(a, b []*Rooms) bool {
	set := make(map[string]bool)
	for _, r := range a[1 : len(a)-1] {
		set[r.Name] = true
	}
	for _, r := range b[1 : len(b)-1] {
		if set[r.Name] {
			return true
		}
	}
	return false
}

// simulateTurns calculates the number of turns needed for a given set of paths.
func simulateTurns(paths [][]*Rooms, totalAnts int) int {
	if len(paths) == 0 {
		return totalAnts + 1 // impossible, return large number
	}

	// sort paths by length
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// distribute ants using the same greedy logic as ShortPath
	nb := make([]int, len(paths))
	remaining := totalAnts
	for remaining > 0 {
		best := 0
		for i := 1; i < len(paths); i++ {
			if len(paths[i])+nb[i] < len(paths[best])+nb[best] {
				best = i
			}
		}
		nb[best]++
		remaining--
	}

	// turns = max(nb[i] + len(paths[i]) - 2) across all paths
	maxTurns := 0
	for i := range paths {
		if nb[i] == 0 {
			continue
		}
		t := nb[i] + len(paths[i]) - 2
		if t > maxTurns {
			maxTurns = t
		}
	}
	return maxTurns
}

// CreatGraph finds all paths via BFS, then selects the best non-overlapping combination.
func CreatGraph() [][]*Rooms {
	allPaths := bfsFindAllPaths()

	if len(allPaths) == 0 {
		log.Fatalln("ERROR: no path between start and end")
	}

	// sort all paths by length
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	totalAnts := len(G.Ants)
	bestTurns := simulateTurns([][]*Rooms{allPaths[0]}, totalAnts)
	bestCombo := [][]*Rooms{allPaths[0]}

	// recursive search for the best non-overlapping combination
	var findBest func(idx int, current [][]*Rooms, usedRooms map[string]bool)
	findBest = func(idx int, current [][]*Rooms, usedRooms map[string]bool) {
		if len(current) > 0 {
			turns := simulateTurns(current, totalAnts)
			if turns < bestTurns {
				bestTurns = turns
				bestCombo = make([][]*Rooms, len(current))
				copy(bestCombo, current)
			}
		}

		for i := idx; i < len(allPaths); i++ {
			// check if this path overlaps with already selected paths
			overlaps := false
			for _, r := range allPaths[i][1 : len(allPaths[i])-1] {
				if usedRooms[r.Name] {
					overlaps = true
					break
				}
			}
			if overlaps {
				continue
			}

			// select this path
			for _, r := range allPaths[i][1 : len(allPaths[i])-1] {
				usedRooms[r.Name] = true
			}
			findBest(i+1, append(current, allPaths[i]), usedRooms)

			// backtrack
			for _, r := range allPaths[i][1 : len(allPaths[i])-1] {
				delete(usedRooms, r.Name)
			}
		}
	}

	findBest(0, nil, make(map[string]bool))
	return bestCombo
}

