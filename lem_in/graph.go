package lem_in

import "fmt"

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

func CreatGraph() [][]*Rooms{
	var all [][]*Rooms
	G.Visited = make(map[string]bool)
	for {
    parent := Bfs(G.RmStar.Name, G.RmEnd.Name)
    if parent == nil {
        break
    }

    path := GeniretPath(parent)
    all = append(all, path)
	if len(path) == 2{
		break
	}
	G.Visited = make(map[string]bool)
    for _, room := range path {
        if room != G.RmStar && room != G.RmEnd {
            G.Visited[room.Name] = true  
        }
    }
}
	fmt.Println("All paths:")
	for _, p := range all {
		for _, r := range p {
			fmt.Print(r.Name, " ")
		}
		fmt.Println()
	}
	return all
}
func Bfs(start string, end string) map[*Rooms]*Rooms {
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

