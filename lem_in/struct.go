package lem_in


type Rooms struct{
	Name string
	X,Y int
	Star bool
	End bool
	IDLinks []int
	IN int
}

type Links struct{
	ID int
	From *Rooms
	To *Rooms
}

type Graph struct{
	Rooms []Rooms
	Links []Links
	Ants []Ants
	CountLinks int
}


type Ants struct{
	ID string
	LocationX int
	LocationY int
	Room *Rooms
	Link *Links
}
func (RM Rooms) Info() (int ,int , bool ,bool){
	return RM.X , RM.Y ,RM.Star ,RM.End
}
func (LK Links) Info() (int ,string , string){
	return LK.ID , LK.From.Name ,LK.To.Name
}

var RmStar *Rooms
var NumberAnts int
var G Graph