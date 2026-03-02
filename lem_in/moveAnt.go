package lem_in

import (
	"fmt"
	"strings"
)


func MoveAnt(all [][]*Rooms ){
	le := len(G.Ants)
	fmt.Println(le)
	all , nb := ShortPath(all)
	res := ""
	for  le != 0{
		for i := 0; i < len(nb); i++ {
			if nb[i] == 0 {
				continue
			}
			for j := len(all[i]) -1 ; j > 0; j-- {
				x := j-1 
				if ((all[i][j].IN == 1 || all[i][j].IN == -2) && (all[i][x].IN == 0 || all[i][x].IN == -2)) && le != 0{
					if all[i][x].IN == -2 {
						if len(G.Ants) == 0{
							fmt.Println(1)
							continue
						}
						fmt.Println(2)
						all[i][j].Ants = &G.Ants[0]
						all[i][j].IN = 1
						G.Ants = G.Ants[1:]
						res += all[i][j].Ants.ID + all[i][j].Name + " "
					}else if (all[i][j].IN == -2){
						
						all[i][x].IN = 0
						all[i][x].Ants = nil
						fmt.Println(res)
						le--
						
					}else {
						fmt.Println(4)
						all[i][j].Ants = all[i][x].Ants
						all[i][x].IN = 0
						all[i][j].IN = 1
						all[i][x].Ants = nil
						res += all[i][j].Ants.ID + all[i][j].Name + " "
					}
				}
			}
		}
		res = strings.TrimSpace(res) + "\n";
	}
	fmt.Println(res)
}

func ShortPath(all [][]*Rooms ) ([][]*Rooms , []int){
	var nb []int
	res := 0
	for i := 0; i < len(all)-1; i++ {
		for j := i+1; j < len(all); j++ {
			if len(all[i]) > len(all[j]) {
				all[i] , all[j] = all[j] , all[i]
			} 
		}
	}
	for i := 0; i < len(all); i++ {
		nb = append(nb, len(all[i]))
		res += len(all[i])
	}
	onePourtion := 100/res
	onePourtionAntDe := len(G.Ants)%100
	onePourtionAnt := (len(G.Ants)-onePourtionAntDe)/100
	for i := 0; i < len(nb); i++ {
		nb[i] = nb[i] * onePourtion
		nb[i] = nb[i] * onePourtionAnt
	}
	nb[0] +=  onePourtionAntDe
	return all , nb
	/*
	10 + 5 + 6 + 4
	100/25   == 4
	10*4 + 5*4 + 6*4 + 4*4
	40 + 20 + 24 + 16
	700/100 == 7
	40*7 + 20*7 + 24*7 + 16*7 -
	280 + 140 + 168 + 112 -
	700 =
	*/
}