package lem_in


func MoveAnt(all [][]*Rooms ){
	le := len(G.Ants)
	all , nb := ShortPath(all)
	for  le != 0{
		for i := 0; i < len(nb); i++ {
			if nb[i] == 0 {
				continue
			}
			for j := len(all[i]) -1 ; j > 0; j-- {
				x := j-1 
				if ((all[i][j].IN == 1 || all[i][j].IN == -2) && (all[i][x].IN == 0 || all[i][x].IN == -2)) && le != 0{
					if all[i][x].IN == -2 {
						l
					}
				}
			}
		}
	}
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
	onePourtionAnt := len(G.Ants)/100
	for i := 0; i < len(nb); i++ {
		nb[i] = nb[i] * onePourtion
		nb[i] = nb[i] * onePourtionAnt
	}
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