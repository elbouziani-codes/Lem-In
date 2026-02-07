package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"Lem-in/lem_in"
)

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
	return "true", true
}

func ValidContent(s string) (string, bool) {
	slice := strings.Split(s, "\n")
	n, err := strconv.Atoi(slice[0])
	if err != nil {
		return err.Error(), false
	}
	if n <= 0 {
		return "Error in Number ants", false
	}
	slice = slice[1:]
	return s, true
}

func validSlice(slice []string, str string){
	countStar := 0
	countEnd := 0
	for i := 0; i < len(slice); i++ {
		if slice[i][0] == '#' {
			if slice[i] == "##start" {
				if countStar == 1 {
					log.Fatalln("Error in content file")
				}
				countStar++
				i += 1
				validRoom(strings.Split(slice[i], " "), 1)
				continue
			} else if slice[i] == "##end" {
				if countEnd == 1 {
					log.Fatalln("Error in content file")
				}
				countEnd++
				i += 1
				validRoom(strings.Split(slice[i], " "), -1)
				continue
			} else {
				continue
			}
		}
		mini_slice := strings.Split(slice[i], " ")
		if len(mini_slice) == 3 {
			validRoom(strings.Split(slice[i], " "), 0)
		} else if len(mini_slice) == 1 && strings.Contains(slice[i],"-"){
			validLink(slice[i])
		} else {
			log.Fatalln("Error in content file")
		}
	}
}

func validRoom(s []string, sore int) {
	if !check(s[0]){
		log.Fatalln("Error in content file")
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
			lem_in.AllRooms = append(lem_in.AllRooms, lem_in.Rooms{
				Name:    s[0],
				X:       n1,
				Y:       n2,
				Star:    true,
				End:     false,
				IDLinks: -1,
				IN:      -1,
			})
		} else if sore == -1 {
			lem_in.AllRooms = append(lem_in.AllRooms, lem_in.Rooms{
				Name:    s[0],
				X:       n1,
				Y:       n2,
				Star:    false,
				End:     true,
				IDLinks: -1,
				IN:      0,
			})
		} else if sore == 0 {
			lem_in.AllRooms = append(lem_in.AllRooms, lem_in.Rooms{
				Name:    s[0],
				X:       n1,
				Y:       n2,
				Star:    false,
				End:     false,
				IDLinks: -1,
				IN:      0,
			})
		}
	} else {
		log.Fatalln("Error in content file")
	}
}
func validLink(s string){
	mini_slice := strings.Split(s,"-")
	count := 0 
	if len(mini_slice) == 2 && mini_slice[0] != mini_slice[1]{
		for i , RM := range lem_in.AllRooms{
			if  RM.Name == mini_slice[0]{
				lem_in.CountLinks++
				lem_in.AllLinks = append(lem_in.AllLinks, lem_in.Links{
					ID:lem_in.CountLinks,
					From:&lem_in.AllRooms[i],
					To:nil,
				})
				count++
				break
			}
		}
		for i , RM := range lem_in.AllRooms{
			if  RM.Name == mini_slice[1]{
				for j , LK := range lem_in.AllLinks{
					if LK.ID ==  lem_in.CountLinks{
						lem_in.AllLinks[j].To = &lem_in.AllRooms[i]
						count++
						break
					}
				}
				break
			}
		}
		if count != 2{
			log.Fatalln("Error in content file")
		}
	}else{
		log.Fatalln("Error in content file")
	}
}

func check(b string) bool{
	for i := 0; i < len(lem_in.AllRooms); i++ {
		if lem_in.AllRooms[i].Name == b {
			return false
		}
	}
	return true
}