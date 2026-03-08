package lem_in

import (
	"log"
)

func GraphRoomsAndLinkes() {
	G.Network = make(map[string][]*Rooms)
	for i, _ := range G.Links {
		G.Links[i].Capacity = 1

		from := G.Links[i].From
		to := G.Links[i].To

		G.Network[from.Name] = append(G.Network[from.Name], to)
		G.Network[to.Name] = append(G.Network[to.Name], from)
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

func Capacity(a, b *Rooms) int {

	for _, link := range G.Links {
		if (link.From.Name == a.Name && link.To.Name == b.Name) || (link.From.Name == b.Name && link.To.Name == a.Name) {
			return link.Capacity
		}
	}
	return 0
}

func CreatGraph() [][]*Rooms {

	multiPath := false
	for _, v := range G.Network {
		if len(v) > 2 {
			multiPath = true
			break
		}
	}
	var all [][]*Rooms
	if multiPath {
		for {
			parent := Bfs(G.RmStar.Name, G.RmEnd.Name)
			if parent == nil {
				if len(all) != 0 {
					return all
				}
				log.Fatalln("Error in path not conection betwine star and end 3")
			}
			path := GeniretPath(parent)
			if path[0] != G.RmStar || path[len(path)-1] != G.RmEnd {
				log.Fatalln("Error in path not conection betwine star and end 2")
			}
			if len(all) > 0 {
				last := all[len(all)-1]
				if equalPath(last, path) {
					break
				}
			}

			all = append(all, path)
			UpdateCapacity(path)

			if len(path) == 2 {
				break
			}
		}
	} else {
		path := Dfs(G.RmStar.Name, G.RmEnd.Name)
		if path == nil {
			log.Fatalln("Error in path not conection betwine star and end ")
		}

		if path[0] != G.RmStar || path[len(path)-1] != G.RmEnd {
			log.Fatalln("Error in path not conection betwine star and end 1")
		}
		all = append(all, path)
	}
	return all
}

func UpdateCapacity(path []*Rooms) {
	for i := 0; i < len(path)-1; i++ {

		a := path[i]
		b := path[i+1]

		for j := range G.Links {

			if G.Links[j].From == a && G.Links[j].To == b {
				G.Links[j].Capacity -= 1
			}

			if G.Links[j].From == b && G.Links[j].To == a {
				G.Links[j].Capacity += 1
			}
		}
	}
}

func Bfs(start string, end string) map[*Rooms]*Rooms {
	G.Visited = make(map[string]bool)
	queue := []*Rooms{}
	queue = append(queue, G.RmStar)
	parent := make(map[*Rooms]*Rooms)
	G.Visited[start] = true
	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]

		for _, next := range G.Network[st.Name] {
			if !G.Visited[next.Name] && Capacity(st, next) > 0 {
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

func Dfs(start string, end string) []*Rooms {
	if G.Visited == nil {
		G.Visited = make(map[string]bool)
	}
	queue := []*Rooms{}
	res := []*Rooms{}
	queue = append(queue, G.RmStar)
	G.Visited[start] = true
	res = append(res, G.RmStar)
	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		for _, next := range G.Network[st.Name] {
			if !G.Visited[next.Name] {
				res = append(res, next)

				if next.Name == end {
					return res
				}

				G.Visited[next.Name] = true
				queue = append(queue, next)
			}
		}
	}
	return []*Rooms{}
}
