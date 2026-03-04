package lem_in

import (
	"fmt"
	"log"
)

func GraphRoomsAndLinkes() {
	G.Network = make(map[string][]*Rooms)
	for _, link := range G.Links {
		G.Network[link.From.Name] = append(G.Network[link.From.Name], link.To)
		G.Network[link.To.Name] = append(G.Network[link.To.Name], link.From)
	}
}

func GeniretPath(parent map[*Rooms]*Rooms) []*Rooms {
	cur := G.RmEnd
	var res []*Rooms
	for cur != nil {
		res = append(res, cur)
		cur = parent[cur]
	}
	// reverse
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func equalPath(a, b []*Rooms) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func CreatGraph() [][]*Rooms {
	var all [][]*Rooms
	usedRooms := make(map[string]bool)

	for {
		G.Visited = make(map[string]bool)

		// منع الغرف الداخلية المستعملة سابقاً
		for name := range usedRooms {
			G.Visited[name] = true
		}
		parent := Bfs(G.RmStar.Name, G.RmEnd.Name)
		if parent == nil {
			if len(all) != 0 {
				return all
			}
			log.Fatalln("Error in path not conection betwine star and end ")
		}

		path := GeniretPath(parent)
		if path[0] == G.RmStar && path[len(path)-1] == G.RmEnd {
			log.Fatalln("Error in path not conection betwine star and end ")
		}
		if len(all) > 0 {
			last := all[len(all)-1]
			if equalPath(last, path) {
				break
			}
		}

		all = append(all, path)

		if len(path) == 2 {
			break
		}
		for _, room := range path {
			if room != G.RmStar && room != G.RmEnd {
				usedRooms[room.Name] = true
			}
		}
	}
	fmt.Println("All paths:")

	for _, p := range all {

		for _, r := range p {
			fmt.Print(r)
		}
		fmt.Println()

	}
	return all
}

func Bfs(start string, end string) map[*Rooms]*Rooms {
	if G.Visited == nil {
		G.Visited = make(map[string]bool)
	}
	queue := []*Rooms{}
	queue = append(queue, G.RmStar)
	parent := make(map[*Rooms]*Rooms)
	G.Visited[start] = true
	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		for _, next := range G.Network[st.Name] {
			if !G.Visited[next.Name] {
				parent[next] = st

				if next.Name == end {
					return parent
				}

				G.Visited[next.Name] = true
				queue = append(queue, next)
			}
		}
	}
	return nil
}
