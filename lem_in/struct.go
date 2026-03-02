package lem_in

type Rooms struct{
	Name string
	X,Y int
	Star bool
	End bool
	IDLinks []int
	IN int
	Ants *Ants
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
	Network map[string][]*Rooms
	Visited map[string]bool
	CountLinks int
	RmStar *Rooms
	RmEnd *Rooms
}

	
type Ants struct{
	ID string
	LocationX int
	LocationY int
	Room *Rooms
}
func (RM Rooms) Info() (int ,int , bool ,bool){
	return RM.X , RM.Y ,RM.Star ,RM.End
}
func (LK Links) Info() (int ,string , string){
	return LK.ID , LK.From.Name ,LK.To.Name
}

var NumberAnts int
var G Graph