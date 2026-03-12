package lem_in

import (
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

// ── Edmonds-Karp max-flow on a node-split graph ──────────────────────

// flowEdge represents a directed edge in the flow network.
type flowEdge struct {
	to, rev int // target node index, index of reverse edge in adj[to]
	cap     int // residual capacity
}

// flowGraph is an adjacency-list flow network.
type flowGraph struct {
	adj [][]flowEdge
	n   int
}

func newFlowGraph(n int) *flowGraph {
	return &flowGraph{adj: make([][]flowEdge, n), n: n}
}

// addEdge adds a directed edge u→v with capacity c (and reverse with 0).
func (fg *flowGraph) addEdge(u, v, c int) {
	fg.adj[u] = append(fg.adj[u], flowEdge{to: v, rev: len(fg.adj[v]), cap: c})
	fg.adj[v] = append(fg.adj[v], flowEdge{to: u, rev: len(fg.adj[u]) - 1, cap: 0})
}

// bfs finds an augmenting path from s to t, returns parent edge info.
func (fg *flowGraph) bfs(s, t int) ([]int, []int) {
	parent := make([]int, fg.n)
	parentEdge := make([]int, fg.n)
	for i := range parent {
		parent[i] = -1
	}
	parent[s] = s
	queue := []int{s}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for i, e := range fg.adj[u] {
			if parent[e.to] == -1 && e.cap > 0 {
				parent[e.to] = u
				parentEdge[e.to] = i
				if e.to == t {
					return parent, parentEdge
				}
				queue = append(queue, e.to)
			}
		}
	}
	return nil, nil
}

// maxFlow runs Edmonds-Karp and returns the max flow value.
func (fg *flowGraph) maxFlow(s, t int) int {
	flow := 0
	for {
		parent, parentEdge := fg.bfs(s, t)
		if parent == nil {
			break
		}
		// find bottleneck
		aug := int(1e9)
		v := t
		for v != s {
			u := parent[v]
			e := &fg.adj[u][parentEdge[v]]
			if e.cap < aug {
				aug = e.cap
			}
			v = u
		}
		// update residuals
		v = t
		for v != s {
			u := parent[v]
			ei := parentEdge[v]
			fg.adj[u][ei].cap -= aug
			fg.adj[v][fg.adj[u][ei].rev].cap += aug
			v = u
		}
		flow += aug
	}
	return flow
}

// ── Path extraction & selection ──────────────────────────────────────

// extractPaths traces flow edges to recover node-disjoint paths in
// the original graph. Each intermediate room was split into in/out;
// we follow the flow from source to sink.
func extractPaths(fg *flowGraph, src, sink int, numRooms int, idxToRoom []*Rooms) [][]*Rooms {
	var paths [][]*Rooms

	for {
		// try to trace one path from src to sink
		path := []*Rooms{idxToRoom[src]} // start room
		cur := src
		visited := make(map[int]bool)
		visited[src] = true
		found := false

		for cur != sink {
			moved := false
			for i := range fg.adj[cur] {
				e := &fg.adj[cur][i]
				// An edge has flow if its reverse edge has cap > 0
				// (originally 0, now increased by flow).
				// For node-split: we care about out→in edges (between different rooms)
				// and in→out edges (within a room).
				revEdge := &fg.adj[e.to][e.rev]
				if revEdge.cap > 0 && !visited[e.to] {
					// consume this unit of flow
					revEdge.cap--
					e.cap++

					nodeIdx := e.to
					// If this is an "in" node (even index), the room is idxToRoom[nodeIdx/2]
					// If this is an "out" node (odd index), the room is idxToRoom[nodeIdx/2]
					roomIdx := nodeIdx
					if numRooms > 0 {
						roomIdx = nodeIdx / 2
					}
					room := idxToRoom[roomIdx]

					// Only add room when we arrive at an "out" node or sink/source
					// For source (idx 0 = in, idx 1 = out): already added
					// For sink: add when we reach it
					// For intermediate: add when we reach _out node (odd index)
					if nodeIdx == sink {
						path = append(path, room)
						found = true
						moved = true
						cur = nodeIdx
						break
					}

					if nodeIdx%2 == 1 { // out-node → means we passed through a room
						path = append(path, room)
					}

					visited[nodeIdx] = true
					cur = nodeIdx
					moved = true
					break
				}
			}
			if !moved {
				break
			}
		}
		if !found {
			break
		}
		paths = append(paths, path)
	}
	return paths
}

// simulateTurns calculates the number of turns needed for a given set of paths.
func simulateTurns(paths [][]*Rooms, totalAnts int) int {
	if len(paths) == 0 {
		return totalAnts + 1
	}

	sorted := make([][]*Rooms, len(paths))
	copy(sorted, paths)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) < len(sorted[j])
	})

	nb := make([]int, len(sorted))
	remaining := totalAnts
	for remaining > 0 {
		best := 0
		for i := 1; i < len(sorted); i++ {
			if len(sorted[i])+nb[i] < len(sorted[best])+nb[best] {
				best = i
			}
		}
		nb[best]++
		remaining--
	}

	maxTurns := 0
	for i := range sorted {
		if nb[i] == 0 {
			continue
		}
		t := nb[i] + len(sorted[i]) - 2
		if t > maxTurns {
			maxTurns = t
		}
	}
	return maxTurns
}

// CreatGraph finds optimal node-disjoint paths using Edmonds-Karp max-flow.
func CreatGraph() [][]*Rooms {
	// Build room name → index mapping
	roomIndex := make(map[string]int)
	rooms := make([]*Rooms, 0)
	for i := range G.Rooms {
		roomIndex[G.Rooms[i].Name] = len(rooms)
		rooms = append(rooms, &G.Rooms[i])
	}

	numRooms := len(rooms)
	startIdx := roomIndex[G.RmStar.Name]
	endIdx := roomIndex[G.RmEnd.Name]

	// Node splitting: each room i becomes two nodes: 2*i (in) and 2*i+1 (out)
	// Source in→out and sink in→out have infinite capacity.
	// All other rooms have in→out capacity 1 (node-disjoint constraint).
	totalNodes := numRooms * 2
	fg := newFlowGraph(totalNodes)

	// idxToRoom maps flow-graph node / 2 → original room pointer
	idxToRoom := make([]*Rooms, numRooms)
	for i, r := range rooms {
		idxToRoom[i] = r
	}

	// Internal edges: in → out
	for i := 0; i < numRooms; i++ {
		if i == startIdx || i == endIdx {
			fg.addEdge(2*i, 2*i+1, numRooms) // effectively infinite
		} else {
			fg.addEdge(2*i, 2*i+1, 1) // node capacity = 1
		}
	}

	// External edges: for each link u-v, add out(u)→in(v) and out(v)→in(u)
	for _, link := range G.Links {
		u := roomIndex[link.From.Name]
		v := roomIndex[link.To.Name]
		fg.addEdge(2*u+1, 2*v, 1)
		fg.addEdge(2*v+1, 2*u, 1)
	}

	src := 2 * startIdx   // source in-node
	sink := 2*endIdx + 1  // sink out-node

	// Run max-flow
	maxFlowVal := fg.maxFlow(src, sink)
	_ = maxFlowVal

	// Extract all flow paths
	allPaths := extractPaths(fg, src, sink, numRooms, idxToRoom)

	if len(allPaths) == 0 {
		// fallback: try simple BFS for a single path
		p := bfsOnePath()
		if p == nil {
			return nil
		}
		return [][]*Rooms{p}
	}

	// Try using 1..len(allPaths) paths, pick the combo with fewest turns
	totalAnts := len(G.Ants)

	// Sort paths by length
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	bestTurns := simulateTurns(allPaths[:1], totalAnts)
	bestCount := 1

	for k := 2; k <= len(allPaths); k++ {
		t := simulateTurns(allPaths[:k], totalAnts)
		if t < bestTurns {
			bestTurns = t
			bestCount = k
		}
	}

	return allPaths[:bestCount]
}

// bfsOnePath finds a single shortest path from start to end (fallback).
func bfsOnePath() []*Rooms {
	parent := make(map[string]*Rooms)
	parent[G.RmStar.Name] = G.RmStar
	queue := []*Rooms{G.RmStar}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == G.RmEnd {
			// reconstruct
			var path []*Rooms
			for c := cur; c != G.RmStar; c = parent[c.Name] {
				path = append([]*Rooms{c}, path...)
			}
			path = append([]*Rooms{G.RmStar}, path...)
			return path
		}

		for _, next := range G.Network[cur.Name] {
			if _, ok := parent[next.Name]; !ok {
				parent[next.Name] = cur
				queue = append(queue, next)
			}
		}
	}
	return nil
}
