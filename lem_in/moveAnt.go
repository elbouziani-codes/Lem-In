package lem_in

import (
	"fmt"
	"strings"
)

func MoveAnt(all [][]*Rooms) {
	all, nb := ShortPath(all)
	res := ""
	move := true
	for move {
		move = false
		line := ""
		for i := 0; i < len(nb) ; i++ {
			for j := len(all[i]) - 1; j > 0 ; j-- {
				x := j - 1
				if all[i][x].IN == 1 && all[i][j].IN == -2{
					move = true
					line += all[i][x].Ants.ID +"-"+ all[i][j].Name + " "
					all[i][x].IN = 0
					all[i][x].Ants = nil
				} else if all[i][x].IN == -2 && all[i][j].IN == 0 && nb[i] > 0{
					nb[i]--
					move = true
					if len(G.Ants) == 0 {
						continue
					}
					all[i][j].Ants = &G.Ants[0]
					all[i][j].IN = 1
					G.Ants = G.Ants[1:]
					line += all[i][j].Ants.ID +"-"+ all[i][j].Name + " "
					fmt.Println(line)
				} else if all[i][x].IN == 1 && all[i][j].IN == 0{
					move = true
					all[i][j].Ants = all[i][x].Ants
					all[i][x].IN = 0
					all[i][j].IN = 1
					all[i][x].Ants = nil
					line += all[i][j].Ants.ID +"-"+ all[i][j].Name + " "
					fmt.Println(line)
				}
			}
		}
		res += strings.TrimSpace(line) + "\n"
		if !move {
			break
		}
	}
	fmt.Println(res)
}

func ShortPath(all [][]*Rooms) ([][]*Rooms, []int) {
	var nb []int

	// ترتيب المسارات من الأقصر للأطول
	for i := 0; i < len(all)-1; i++ {
		for j := i + 1; j < len(all); j++ {
			if len(all[i]) > len(all[j]) {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	totalAnts := len(G.Ants)

	for i := 0; i < len(all); i++ {
		capacity := totalAnts - (len(all[i]) - 1)
		if capacity < 0 {
			capacity = 0
		}
		nb = append(nb, capacity)
	}

	return all, nb
}
