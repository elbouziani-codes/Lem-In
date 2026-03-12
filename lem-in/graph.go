package main

// Room represents a single room in the ant colony graph
type Room struct {
	Name string
	X    int
	Y    int
}

// Graph represents the entire ant colony map
type Graph struct {
	NumAnts int
	Start   *Room
	End     *Room
	Rooms   map[string]*Room
	Edges   map[string][]string // adjacency list maps a Room's Name to an array of linked Room Names
}

// NewGraph initializes an empty Graph
func NewGraph() *Graph {
	return &Graph{
		Rooms: make(map[string]*Room),
		Edges: make(map[string][]string),
	}
}

// AddRoom safely adds a room to the graph, returning false if it already exists
func (g *Graph) AddRoom(r *Room) bool {
	if _, exists := g.Rooms[r.Name]; exists {
		return false // duplicate room
	}
	// Note: According to the requirements, duplicate coordinates could also be considered invalid
	for _, existingRoom := range g.Rooms {
		if existingRoom.X == r.X && existingRoom.Y == r.Y {
			return false // coordinates already used
		}
	}
	g.Rooms[r.Name] = r
	return true
}

// AddLink creates a bidirectional path between room1 and room2.
func (g *Graph) AddLink(room1, room2 string) error {
	if room1 == room2 {
		return ErrRoomLinkToItself
	}
	if _, ok := g.Rooms[room1]; !ok {
		return ErrLinkToUnknownRoom
	}
	if _, ok := g.Rooms[room2]; !ok {
		return ErrLinkToUnknownRoom
	}

	// Check for duplicate link
	for _, neighbor := range g.Edges[room1] {
		if neighbor == room2 {
			return ErrDuplicateLink
		}
	}

	g.Edges[room1] = append(g.Edges[room1], room2)
	g.Edges[room2] = append(g.Edges[room2], room1)

	return nil
}
