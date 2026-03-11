package lem_in

import (
	"fmt"
	"strings"
)

const (
	EMPTY = 0
	Full   = 1
	START = -1
	END   = -2
)


func MoveAnt(all [][]*Rooms) {
	all, nb := ShortPath(all)
	res := ""
	move := true
	for move {
		move = false
		line := ""
		for i := 0; i < len(nb); i++ {
			for j := len(all[i]) - 1; j > 0; j-- {
				x := j - 1
				if all[i][x].IN == Full && all[i][j].IN == END {

					line += all[i][x].Ants.ID + "-" + all[i][j].Name + " "
					all[i][x].IN = 0
					all[i][x].Ants = nil
					move = true
				} else if all[i][x].IN == START && all[i][j].IN == EMPTY && nb[i] > 0 {
					nb[i]--
					move = true
					if len(G.Ants) == 0 {
						continue
					}
					all[i][j].Ants = &G.Ants[0]
					all[i][j].IN = 1
					G.Ants = G.Ants[1:]
					line += all[i][j].Ants.ID + "-" + all[i][j].Name + " "
				} else if all[i][x].IN == Full && all[i][j].IN == EMPTY {
					move = true
					all[i][j].Ants = all[i][x].Ants
					all[i][x].IN = 0
					all[i][j].IN = 1
					all[i][x].Ants = nil
					line += all[i][j].Ants.ID + "-" + all[i][j].Name + " "
				} else if all[i][x].IN == START && all[i][j].IN == END && nb[i] > 0 {
					nb[i]--
					move = true
					if len(G.Ants) == 0 {
						continue
					}
					line += G.Ants[0].ID + "-" + all[i][j].Name + " "
					G.Ants = G.Ants[1:]
				}
			}
		}
		if line != "" {
			res += strings.TrimSpace(line) + "\n"
		}
		if !move {
			break
		}
	}
	fmt.Println(res)
}

func ShortPath(all [][]*Rooms) ([][]*Rooms, []int) {

	// ترتيب المسارات من الأقصر للأطول
	for i := 0; i < len(all)-1; i++ {
		for j := i + 1; j < len(all); j++ {
			if len(all[i]) > len(all[j]) {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	nb := make([]int, len(all))
	totalAnts := len(G.Ants)

	for totalAnts > 0 {
		best := 0
		for i := 1; i < len(all); i++ {
			if len(all[i])+nb[i] < len(all[best])+nb[best] {
				best = i
			}
		}
		nb[best]++
		totalAnts--
	}

	return all, nb
}
