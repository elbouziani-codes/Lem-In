package lem_in

type Rooms struct{
	Name string
	X,Y int
	Star bool
	End bool
	IDLinks int
	IN int
}

type Links struct{
	ID int
	From *Rooms
	To *Rooms
}

type Ants struct{
	ID int
	LocationX,LocationY int
	Room Rooms
	Link Links
}

// func (rm Rooms) check() bool{
// 	for i := 0; i < len(AllRooms); i++ {
// 			if AllRooms[i].Name == rm.Name {
// 				return true
// 			}
// 	}
// 	return false
// }
var AllRooms []Rooms
var AllLinks []Links
var AllAnts []Ants
var CountLinks int = 0