package main

// REdge represents a directed edge in the residual graph.
type REdge struct {
	To       int
	Capacity int
	Flow     int
	Cost     int
	Rev      int // Index of the reverse edge in the adjacency list of 'To'
}

// FindOptimalPaths returns a slice of disjoint paths (each path is a slice of Rooms, not including Start and End).
// It models Bhandari's algorithm to find vertex-disjoint shortest paths.
func FindOptimalPaths(graph *Graph) ([][]*Room, error) {
	rooms := make([]*Room, len(graph.Rooms))
	roomIdx := make(map[string]int)

	rooms[0] = graph.Start
	roomIdx[graph.Start.Name] = 0

	rooms[1] = graph.End
	roomIdx[graph.End.Name] = 1

	idx := 2
	for _, r := range graph.Rooms {
		if r == graph.Start || r == graph.End {
			continue
		}
		rooms[idx] = r
		roomIdx[r.Name] = idx
		idx++
	}

	numNodes := 2 * len(rooms)
	adj := make([][]*REdge, numNodes)

	addEdge := func(u, v, cap, cost int) {
		adj[u] = append(adj[u], &REdge{To: v, Capacity: cap, Cost: cost, Rev: len(adj[v])})
		adj[v] = append(adj[v], &REdge{To: u, Capacity: 0, Cost: -cost, Rev: len(adj[u]) - 1})
	}

	// Add vertex capacities for intermediate rooms
	for i := 2; i < len(rooms); i++ {
		inNode := 2 * i
		outNode := 2*i + 1
		addEdge(inNode, outNode, 1, 0)
	}

	for uName, neighbors := range graph.Edges {
		u := roomIdx[uName]
		for _, vName := range neighbors {
			v := roomIdx[vName]
			// Don't add edges going INTO Start or OUT OF End
			if v == 0 || u == 1 {
				continue
			}
			uOut := 2*u + 1
			vIn := 2 * v
			addEdge(uOut, vIn, 1, 1) // Cost 1 implies passing a tunnel increases length by 1
		}
	}

	var bestPaths [][]int
	bestTurns := -1

	// Successive Shortest Path (using SPFA)
	for {
		dist, parentEdge := spfa(adj, 1, 2) // Start_out is 1, End_in is 2
		if dist == nil {
			break
		}

		// Augment flow along the shortest path
		curr := 2
		for curr != 1 {
			edge := parentEdge[curr]
			edge.Flow += 1
			adj[edge.To][edge.Rev].Flow -= 1
			curr = adj[edge.To][edge.Rev].To
		}

		// Extract disjoint paths from the residual graph
		paths := extractPaths(adj)

		// Calculate turns required using the current suite of paths
		turns := calcTurns(paths, graph.NumAnts)

		// Wait, if paths is 0 length, it means no paths could be extracted (shouldn't happen here)
		if len(paths) == 0 {
			break
		}

		// Keep the set of paths that produces the absolute minimum number of turns
		if bestTurns == -1 || turns < bestTurns {
			bestTurns = turns
			bestPaths = make([][]int, len(paths))
			for i, p := range paths {
				bestPaths[i] = make([]int, len(p))
				copy(bestPaths[i], p)
			}
		} else {
			// Because length is strictly monotonic and turns form a convex curve,
			// once the number of turns increases, we have passed the optimal set of paths.
			break
		}
	}

	if bestTurns == -1 {
		return nil, ErrNoPath
	}

	// Sort the bestPaths by length from shortest to longest for simulation
	for i := 0; i < len(bestPaths); i++ {
		for j := i + 1; j < len(bestPaths); j++ {
			if len(bestPaths[i]) > len(bestPaths[j]) {
				bestPaths[i], bestPaths[j] = bestPaths[j], bestPaths[i]
			}
		}
	}

	var result [][]*Room
	for _, p := range bestPaths {
		var pathRooms []*Room
		// p contains room indices
		for _, rIdx := range p {
			pathRooms = append(pathRooms, rooms[rIdx])
		}
		result = append(result, pathRooms)
	}

	return result, nil
}

// spfa is the Shortest Path Faster Algorithm capable of handling negative edge weights.
func spfa(adj [][]*REdge, s, t int) ([]int, []*REdge) {
	n := len(adj)
	dist := make([]int, n)
	parentEdge := make([]*REdge, n)
	inQueue := make([]bool, n)
	popCount := make([]int, n)

	for i := range dist {
		dist[i] = 1e9 // equivalent to infinity
	}

	queue := make([]int, 0, n)
	queue = append(queue, s)
	dist[s] = 0
	inQueue[s] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false
		popCount[u]++
		
		// Unlikely to hit negative cycles as residual graph of EK has none, but added for safety
		if popCount[u] > n {
			return nil, nil
		}

		for _, edge := range adj[u] {
			if edge.Capacity-edge.Flow > 0 { // Can push flow
				if dist[u]+edge.Cost < dist[edge.To] {
					dist[edge.To] = dist[u] + edge.Cost
					parentEdge[edge.To] = edge
					if !inQueue[edge.To] {
						queue = append(queue, edge.To)
						inQueue[edge.To] = true
					}
				}
			}
		}
	}

	if dist[t] == 1e9 {
		return nil, nil // No path
	}
	return dist, parentEdge
}

// extractPaths follows the original forward edges with Flow == 1 to build separate paths.
func extractPaths(adj [][]*REdge) [][]int {
	var paths [][]int
	startOut := 1
	endIn := 2

	for _, edge := range adj[startOut] {
		if edge.Capacity > 0 && edge.Flow == 1 {
			path := []int{}
			currNode := edge.To
			for currNode != endIn {
				roomIdx := currNode / 2
				path = append(path, roomIdx)

				outNode := 2*roomIdx + 1
				moved := false
				for _, nextEdge := range adj[outNode] {
					if nextEdge.Capacity > 0 && nextEdge.Flow == 1 {
						currNode = nextEdge.To
						moved = true
						break
					}
				}
				if !moved {
					break
				}
			}
			paths = append(paths, path)
		}
	}
	return paths
}

// calcTurns determines the minimum number of turns needed to send numAnts through given paths.
func calcTurns(paths [][]int, numAnts int) int {
	if len(paths) == 0 {
		return 1e9
	}
	// Calculate T directly using binary search
	low := 1
	high := numAnts + 1000000 // Upper bound
	bestT := high

	for low <= high {
		mid := (low + high) / 2

		capacity := 0
		for _, p := range paths {
			l := len(p) + 1 // Length in edges
			if mid >= l {
				capacity += mid - l + 1
			}
		}

		if capacity >= numAnts {
			bestT = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return bestT
}
