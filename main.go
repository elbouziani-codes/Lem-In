package main

import (
	"log"
	"os"
)
func main(){
	if len(os.Args) != 2 {
		log.Fatalln("Error in length arguments ")
	}
	Content , BOOL :=ReadAllLines(os.Args[1])
	if !BOOL{
		log.Fatalln(Content)
	}
}
	
