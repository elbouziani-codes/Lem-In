package main

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"Lem-in/lem_in"
)

var xStart, yStart int

func ReadAllLines(NameFile string) (string, bool) {
	newName := strings.ToLower(NameFile)
	if !strings.HasSuffix(newName, ".txt") {
		return "Error in file prefix: example NameFile.txt", false
	}
	content, err := os.ReadFile(NameFile)
	if err != nil {
		return err.Error(), false
	}
	str, BOOL := ValidContent(string(content))
	if !BOOL {
		return str, false
	}

	validRoomsRepet()
	validLinksRepet()
	parsingAnts()
	checkRoomLinkes()
	lem_in.GraphRoomsAndLinkes()
	lem_in.MoveAnt(lem_in.CreatGraph())
	return str, true
}

func ValidContent(s string) (string, bool) {
	slice := strings.Split(s, "\n")
	i := 0
	for ; i < len(slice); i++ {
		if len(slice[i]) == 0 {
			continue
		} else {
			numberAnts, err := strconv.Atoi(slice[i])
			if err != nil {
				return err.Error(), false
			}
			if numberAnts <= 0 {
				return "Error in Number ants", false
			}
			lem_in.NumberAnts = numberAnts
			break
		}
	}
	slice = slice[i+1:]
	validSlice(slice)
	return s, true
}

func validSlice(slice []string) {
	countStar := 0
	countEnd := 0
	i := 0
	for ; i < len(slice); i++ {
		if len(slice[i]) == 0 {
			continue
		}
		if slice[i][0] == '#' {
			if slice[i] == "##start" {
				if countStar == 1 {
					log.Fatalln("Error in content file 75")
				}
				countStar++
				i += 1
				if i >= len(slice) {
					log.Fatal("Error in content file 80")
				}
				validRoom(strings.Split(slice[i], " "), 1)
				continue
			} else if slice[i] == "##end" {
				if countEnd == 1 {
					log.Fatalln("Error in content file 86")
				}
				countEnd++
				i += 1
				if i >= len(slice) {
					log.Fatal("Error in content file 91")
				}
				validRoom(strings.Split(slice[i], " "), -1)
				continue
			} else {
				continue
			}
		}
		mini_slice := strings.Split(strings.TrimSpace(slice[i]), " ")
		if len(mini_slice) == 3 {
			validRoom(strings.Split(strings.TrimSpace(slice[i]), " "), 0)
		} else {
			break
		}
	}
	for ; i < len(slice); i++ {
		if len(slice[i]) == 0 {
			continue
		} else if slice[i][0] == '#' {
			continue
		} else if strings.Contains(slice[i], "-") {
			validLink(strings.TrimSpace(slice[i]))
		} else {
			log.Fatalln("Error in content file 112")
		}
	}
}

func validRoom(s []string, sore int) {
	if !check(s[0]) {
		log.Fatalln("Error in content file 119")
	}
	if len(s) == 3 {
		n1, err := strconv.Atoi(s[1])
		if err != nil {
			log.Fatalln(err.Error())
		}
		n2, err := strconv.Atoi(s[2])
		if err != nil {
			log.Fatalln(err.Error())
		}
		if sore == 1 {
			lem_in.G.Rooms = append(lem_in.G.Rooms, lem_in.Rooms{
				Name:    s[0],
				X:       n1,
				Y:       n2,
				Star:    true,
				End:     false,
				IDLinks: []int{},
				IN:      lem_in.START,
				Ants:    nil,
			})
			lem_in.G.RmStar = &lem_in.G.Rooms[len(lem_in.G.Rooms)-1]
			xStart = n1
			yStart = n2

		} else if sore == -1 {
			lem_in.G.Rooms = append(lem_in.G.Rooms, lem_in.Rooms{
				Name:    s[0],
				X:       n1,
				Y:       n2,
				Star:    false,
				End:     true,
				IDLinks: []int{},
				IN:      lem_in.END,
				Ants:    nil,
			})
			lem_in.G.RmEnd = &lem_in.G.Rooms[len(lem_in.G.Rooms)-1]

		} else if sore == 0 {
			lem_in.G.Rooms = append(lem_in.G.Rooms, lem_in.Rooms{
				Name:    s[0],
				X:       n1,
				Y:       n2,
				Star:    false,
				End:     false,
				IDLinks: []int{},
				IN:      lem_in.EMPTY,
				Ants:    nil,
			})
		}
		// room.IN == -2 => capacity all ant -2 in start and end
		// room.IN == 0 =>  0 ant max capacity =1
		// room.IN == 1 => 1 ant full capacity
	} else {
		log.Fatalln("Error in content file 171")
	}
}

func validLink(s string) {
	mini_slice := strings.Split(s, "-")
	count := 0
	if len(mini_slice) == 2 && strings.TrimSpace(mini_slice[0]) != strings.TrimSpace(mini_slice[1]) {
		for i, RM := range lem_in.G.Rooms {
			if RM.Name == strings.TrimSpace(mini_slice[0]) {
				lem_in.G.CountLinks++
				lem_in.G.Links = append(lem_in.G.Links, lem_in.Links{
					ID:   lem_in.G.CountLinks,
					From: &lem_in.G.Rooms[i],
					To:   nil,
					Capacity: 1,
				})
				count++
				break
			}
		}
		for i, RM := range lem_in.G.Rooms {
			if RM.Name == strings.TrimSpace(mini_slice[1]) {
				for j, LK := range lem_in.G.Links {
					if LK.ID == lem_in.G.CountLinks {
						lem_in.G.Links[j].To = &lem_in.G.Rooms[i]
						count++
						break
					}
				}
				break
			}
		}
		if count != 2 {
			log.Fatalln("Error in content file 203")
		}
	} else {
		log.Fatalln("Error in content file 206")
	}
}

func check(b string) bool {
	for i := 0; i < len(lem_in.G.Rooms); i++ {
		if lem_in.G.Rooms[i].Name == b {
			return false
		}
	}
	return true
}

// validate room is repeat
func validRoomsRepet() {
	for i := 0; i < len(lem_in.G.Rooms)-1; i++ {
		for j := i + 1; j < len(lem_in.G.Rooms); j++ {
			X1, Y1, _, _ := lem_in.G.Rooms[i].Info()
			X2, Y2, _, _ := lem_in.G.Rooms[j].Info()
			if (X1 == X2) && (Y1 == Y2) {
				log.Fatal("Erorr in content File 225")
			}
		}
	}
	star := false
	end := false
	for i := 0; i < len(lem_in.G.Rooms); i++ {
		_, _, s, e := lem_in.G.Rooms[i].Info()
		if s && e {
			log.Fatal("Erorr in content File 235")
		} else if s && !e {
			star = true
		} else if e && !s {
			end = true
		} else {
			continue
		}
		if star && end {
		}
	}
	if !star || !end {
		log.Fatal("Erorr in content File 248")
	}
}

func validLinksRepet() {
	for i := 0; i < len(lem_in.G.Links)-1; i++ {
		for j := i + 1; j < len(lem_in.G.Links); j++ {
			ID1, RM11, RM12 := lem_in.G.Links[i].Info()
			ID2, RM21, RM22 := lem_in.G.Links[j].Info()
			if ((RM11 == RM21) && (RM12 == RM22)) || ID1 == ID2 {
				log.Fatal("Erorr in content File 258")
			} else if (RM11 == RM22) && (RM12 == RM21) {
				log.Fatal("Erorr in content File 260")
			}

		}
	}
	for i := 0; i < len(lem_in.G.Links); i++ {
		ID, rm1, rm2 := lem_in.G.Links[i].Info()
		if rm1 == rm2 {
			log.Fatal("Erorr in content File 269")
		}
		if !slices.Contains(lem_in.G.Links[i].From.IDLinks, ID) {
			lem_in.G.Links[i].From.IDLinks = append(lem_in.G.Links[i].From.IDLinks, ID)
		}
		if !slices.Contains(lem_in.G.Links[i].To.IDLinks, ID) {
			lem_in.G.Links[i].To.IDLinks = append(lem_in.G.Links[i].To.IDLinks, ID)
		}
	}
}

func parsingAnts() {
	for i := 1; i <= lem_in.NumberAnts; i++ {
		lem_in.G.Ants = append(lem_in.G.Ants, lem_in.Ants{
			ID:        "L" + strconv.Itoa(i),
			LocationX: xStart,
			LocationY: yStart,
			Room:      lem_in.G.RmStar,
		})
	}
}

func checkRoomLinkes() {
	for i := len(lem_in.G.Rooms) - 1; i >= 0; i-- {
		if len(lem_in.G.Rooms[i].IDLinks) == 0 {
			if lem_in.G.Rooms[i].Star || lem_in.G.Rooms[i].End {
				log.Fatalln("Error in content File 297")
			} else {
				lem_in.G.Rooms = append(lem_in.G.Rooms[:i], lem_in.G.Rooms[i+1:]...)
			}
		}
	}
}
